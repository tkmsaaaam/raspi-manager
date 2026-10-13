package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

const accept = "application/vnd.oci.image.index.v1+json, application/vnd.docker.distribution.manifest.list.v2+json, application/vnd.oci.image.manifest.v1+json, application/vnd.docker.distribution.manifest.v2+json"

type entry struct {
	Name, Digest string
	Tags         []string
	Created      *time.Time
}
type manifest struct {
	MediaType string `json:"mediaType"`
	Config    struct {
		Digest string `json:"digest"`
	} `json:"config"`
	Manifests []struct {
		Digest   string                            `json:"digest"`
		Platform struct{ OS, Architecture string } `json:"platform"`
	} `json:"manifests"`
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s http://localhost:5000", os.Args[0])
	}
	base := strings.TrimRight(os.Args[1], "/")
	logger := log.New(os.Stdout, "", 0)
	client := &http.Client{Timeout: 30 * time.Second}
	get := func(raw string, acceptHeader string) ([]byte, http.Header, error) {
		req, err := http.NewRequest(http.MethodGet, raw, nil)
		if err != nil {
			return nil, nil, err
		}
		if acceptHeader != "" {
			req.Header.Set("Accept", acceptHeader)
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return nil, nil, fmt.Errorf("GET %s: HTTP %d", raw, resp.StatusCode)
		}
		b, err := io.ReadAll(resp.Body)
		return b, resp.Header, err
	}
	list := func(start, field string) ([]string, error) {
		var out []string
		for next := start; next != ""; {
			b, h, err := get(next, "")
			if err != nil {
				return nil, err
			}
			var v struct {
				Repositories []string `json:"repositories"`
				Tags         []string `json:"tags"`
			}
			if err = json.Unmarshal(b, &v); err != nil {
				return nil, err
			}
			if field == "repositories" {
				out = append(out, v.Repositories...)
			} else {
				out = append(out, v.Tags...)
			}
			old := next
			next = ""
			for _, part := range strings.Split(h.Get("Link"), ",") {
				part = strings.TrimSpace(part)
				if !strings.Contains(part, "rel=\"next\"") && !strings.Contains(part, "rel=next") {
					continue
				}
				l, r := strings.Index(part, "<"), strings.Index(part, ">")
				if l < 0 || r <= l {
					continue
				}
				ref, e := url.Parse(part[l+1 : r])
				if e != nil {
					return nil, e
				}
				u, e := url.Parse(old)
				if e != nil {
					return nil, e
				}
				next = u.ResolveReference(ref).String()
				break
			}
		}
		return out, nil
	}
	var resolve func(string, string, int) (string, *time.Time, error)
	resolve = func(repo, ref string, depth int) (string, *time.Time, error) {
		if depth > 4 {
			return "", nil, fmt.Errorf("manifest nesting too deep")
		}
		b, h, err := get(base+"/v2/"+repo+"/manifests/"+url.PathEscape(ref), accept)
		if err != nil {
			return "", nil, err
		}
		digest := h.Get("Docker-Content-Digest")
		if digest == "" {
			return "", nil, fmt.Errorf("missing Docker-Content-Digest")
		}
		var m manifest
		if err = json.Unmarshal(b, &m); err != nil {
			return "", nil, err
		}
		if len(m.Manifests) > 0 {
			// Prefer a runnable platform over an attestation manifest.
			for _, child := range m.Manifests {
				if child.Platform.OS == "linux" && child.Platform.Architecture == "amd64" {
					_, created, e := resolve(repo, child.Digest, depth+1)
					return digest, created, e
				}
			}
			for _, child := range m.Manifests {
				if child.Platform.OS != "unknown" && child.Platform.OS != "" {
					_, created, e := resolve(repo, child.Digest, depth+1)
					return digest, created, e
				}
			}
			return digest, nil, nil
		}
		if m.Config.Digest == "" {
			return digest, nil, nil
		}
		blob, _, err := get(base+"/v2/"+repo+"/blobs/"+url.PathEscape(m.Config.Digest), "")
		if err != nil {
			return digest, nil, err
		}
		var cfg struct {
			Created *time.Time `json:"created"`
		}
		if err = json.Unmarshal(blob, &cfg); err != nil {
			return digest, nil, err
		}
		return digest, cfg.Created, nil
	}
	repos, err := list(base+"/v2/_catalog?n=100", "repositories")
	if err != nil {
		log.Fatal(err)
	}
	var entries []entry
	byDigest := make(map[string]int)
	for _, repo := range repos {
		tags, e := list(base+"/v2/"+repo+"/tags/list?n=100", "tags")
		if e != nil {
			log.Printf("%s: %v", repo, e)
			continue
		}
		for _, tag := range tags {
			digest, created, e := resolve(repo, tag, 0)
			if e != nil {
				log.Printf("%s:%s: %v", repo, tag, e)
				continue
			}
			key := repo + "\x00" + digest
			if i, ok := byDigest[key]; ok {
				entries[i].Tags = append(entries[i].Tags, tag)
				if entries[i].Created == nil {
					entries[i].Created = created
				}
			} else {
				byDigest[key] = len(entries)
				entries = append(entries, entry{Name: repo, Digest: digest, Tags: []string{tag}, Created: created})
			}
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		a, b := entries[i], entries[j]
		if a.Created == nil {
			return false
		}
		if b.Created == nil {
			return true
		}
		if !a.Created.Equal(*b.Created) {
			return a.Created.Before(*b.Created)
		}
		if a.Name != b.Name {
			return a.Name < b.Name
		}
		return a.Digest < b.Digest
	})
	for _, e := range entries {
		sort.Strings(e.Tags)
		created := "unknown"
		if e.Created != nil {
			created = e.Created.UTC().Format(time.RFC3339)
		}
		deleteURL := base + "/v2/" + e.Name + "/manifests/" + e.Digest
		deleteURL = strings.ReplaceAll(deleteURL, "'", "'\"'\"'")
		logger.Printf("created=%s pushed=unknown %s:%s -> %s | curl -X DELETE '%s'", created, e.Name, "["+strings.Join(e.Tags, ",")+"]", e.Digest, deleteURL)
	}

}
