[<EntryPoint>]
let main args =
    let today = System.DateTime.Today
    let filterMonth (month: string) =
        month.Equals(today.ToString "MMM" )
    let filterDate (date: string) =
        date.Equals(string today.Day)
    let filterDaemon (daemon: string) =
        daemon.StartsWith("sshd") || daemon.StartsWith("systemd-logind")
    let filterLine (line: string) = // TODO: fix filterDate, filterDaemon
        line.Split()[0] |> filterMonth && line.Split()[2] |> filterDate && line.Split()[5] |> filterDaemon
    let file = @"/logs/auth.log"
    let lines = System.IO.File.ReadAllLines(file)
    let filterLines (lines: string array) =
        lines |> Array.filter (fun line -> filterLine line)
    for line in filterLines lines do
        printfn "%s" line
    0
