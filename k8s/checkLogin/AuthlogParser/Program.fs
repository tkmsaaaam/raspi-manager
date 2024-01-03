[<EntryPoint>]
let main args =
    let today = System.DateTime.Today
    let filterDate (date: string) =
        date.Equals(today.ToString "yyyy/MM/dd" )
    let filterDaemon (daemon: string) =
        daemon.StartsWith("sshd") || daemon.StartsWith("systemd-logind")
    let filterLine (line: string) =
        line.Split("|")[0] |> filterDate && line.Split("|")[5] |> filterDaemon
    let file = @"/logs/authlog"
    let lines = System.IO.File.ReadAllLines(file)
    let filterLines (lines: string array) =
        lines |> Array.filter (fun line -> filterLine line)
    for line in filterLines lines do
        printfn "%s" line
    0
