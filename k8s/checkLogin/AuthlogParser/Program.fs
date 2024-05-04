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
    let timeEnd = 15
    let hostStart = timeEnd + 1
    let hostname = System.Environment.GetEnvironmentVariable("HOSTNAME")
    let hostnameLength = hostname.Length
    let daemonStart = hostStart + hostnameLength + 1

    let filterDate (date: string) =
        date.StartsWith(today.ToString("MMM", format))
        && date[3..5].EndsWith(today.Day.ToString())

    let filterDaemon (daemon: string) =
        daemon.StartsWith("sshd") || daemon.StartsWith("systemd-logind")

    let filterLine (line: string) =
        line |> filterDate && line[daemonStart..] |> filterDaemon

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
                      line[0..timeEnd],
                      line[hostStart .. hostStart + hostnameLength],
                      line[daemonStart .. i - 1],
                      line[i + 2 ..]
                  ) ]


    System.Text.Json.JsonSerializer.Serialize<Result>(Result(today.ToString("MMM dd", format), hostname, list))
    |> printf "%s"

    0
