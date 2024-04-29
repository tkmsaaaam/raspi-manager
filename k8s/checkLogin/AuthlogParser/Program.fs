[<EntryPoint>]
let main args =
    let jstTimeZone =
      System.TimeZoneInfo.FindSystemTimeZoneById("Tokyo Standard Time");
    let today =
      System.TimeZoneInfo.ConvertTime(System.DateTime.Today, jstTimeZone);
    let timeLength =
      today.ToString("MMM d HH:mm:ss", System.Globalization.CultureInfo.CreateSpecificCulture("en-US")).Length
    let hostnameLength =  System.Environment.GetEnvironmentVariable("HOSTNAME").Length
    let filterDate (date: string) =
      date.StartsWith(today.ToString("MMM d", System.Globalization.CultureInfo.CreateSpecificCulture("en-US")))
    let filterDaemon (daemon: string) =
      daemon.StartsWith("sshd") || daemon.StartsWith("systemd-logind")
    let filterLine (line: string) =
      line |> filterDate && line[timeLength+1+hostnameLength+1..] |> filterDaemon
    let file = @"/logs/auth.log"
    let lines = System.IO.File.ReadAllLines(file)
    let filterLines (lines: string array) =
      lines |> Array.filter (fun line -> filterLine line)
    for line in filterLines lines do
      line |> printfn "%s"
    0
