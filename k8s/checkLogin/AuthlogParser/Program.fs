type Log =
    struct
        val dateTime: string
        val host: string
        val daemon: string
        val message: string

        new(dateTime: string, host: string, daemon: string, message: string) =
            { dateTime = dateTime
              host = host
              daemon = daemon
              message = message }
    end

type Result =
    struct
        val date: string
        val host: string
        val logs: Log list

        new(date: string, host: string, logs: Log list) =
            { date = date
              host = host
              logs = logs }
    end

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
        && line[timeLength + 1 + hostnameLength + 2 ..] |> filterDaemon

    let file = @"/logs/auth.log"
    let lines = System.IO.File.ReadAllLines(file)

    let filterLines (lines: string array) =
        lines |> Array.filter (fun line -> filterLine line)

    let mutable list = []


    for line in filterLines lines do
        let i = line.IndexOf(": ")

        list <-
            List.append
                list
                [ Log(
                      line[0..timeLength],
                      line[timeLength + 1 .. timeLength + 1 + hostnameLength],
                      line[timeLength + 1 + hostnameLength + 2 .. i - 1],
                      line[i + 2 ..]
                  ) ]


    System.Text.Json.JsonSerializer.Serialize<Result>(Result(todayMmmD, hostname, list))
    |> printf "%s"

    0
