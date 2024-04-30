[<EntryPoint>]
let main args =
    let jstTimeZone = System.TimeZoneInfo.FindSystemTimeZoneById("Tokyo Standard Time")
    let format = System.Globalization.CultureInfo.CreateSpecificCulture("en-US")
    let today = System.TimeZoneInfo.ConvertTime(System.DateTime.Today, jstTimeZone)
    let todayMmmD = today.ToString("MMM d", format)
    let timeLength = today.ToString("MMM d HH:mm:ss", format).Length
    let hostname = System.Environment.GetEnvironmentVariable("HOSTNAME")
    let hostnameLength = hostname.Length
    let filterDate (date: string) = date.StartsWith(todayMmmD)

    let filterDaemon (daemon: string) =
        daemon.StartsWith("sshd") || daemon.StartsWith("systemd-logind")

    let filterLine (line: string) =
        line |> filterDate
        && line[timeLength + 1 + hostnameLength + 1 ..] |> filterDaemon

    let file = @"/logs/auth.log"
    let lines = System.IO.File.ReadAllLines(file)

    let filterLines (lines: string array) =
        lines |> Array.filter (fun line -> filterLine line)

    printfn "searching..."
    todayMmmD |> printfn "date: %s"
    hostname |> printfn "host: %s"

    for line in filterLines lines do
        line |> printfn "%s"

    0
