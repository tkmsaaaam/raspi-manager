module Tests

open Xunit

open Program

let jstTimeZone = System.TimeZoneInfo.FindSystemTimeZoneById("Tokyo Standard Time")
let format = System.Globalization.CultureInfo.CreateSpecificCulture("en-US")

let dateTime = new System.DateTime(2001, 1, 1)
let today = System.TimeZoneInfo.ConvertTime(dateTime, jstTimeZone)

[<Fact>]
let `` Log#new is Ok `` () =
    let dateTime = "2000-01-01 00:00:00"
    let host = "hostname"
    let daemon = "daemon name"
    let message = "message detail"
    let actual = Log(dateTime, host, daemon, message)
    Assert.Equal(dateTime, actual.dateTime)
    Assert.Equal(host, actual.host)
    Assert.Equal(daemon, actual.daemon)
    Assert.Equal(message, actual.message)

[<Fact>]
let ``Result#new is Ok`` () =
    let date = "May  1"
    let host = "hostname"
    let log1 = Log("2000-01-01 00:00:00", "hostname", "daemon", "message detail")
    let log2 = Log("2000-01-02 00:00:00", "hostname2", "daemon2", "message detail")
    let logs = List.append [ log1 ] [ log2 ]
    let actual = Result(date, host, logs)
    Assert.Equal(date, actual.date)
    Assert.Equal(host, actual.host)
    Assert.Equal(logs.Length, actual.logs.Length)
    Assert.Equal(log1.dateTime, actual.logs[0].dateTime)
    Assert.Equal(log1.host, actual.logs[0].host)
    Assert.Equal(log1.daemon, actual.logs[0].daemon)
    Assert.Equal(log2.dateTime, actual.logs[1].dateTime)
    Assert.Equal(log2.host, actual.logs[1].host)
    Assert.Equal(log2.daemon, actual.logs[1].daemon)

[<Fact>]
let ``filterDate is Ok`` () =
    let line = "Jan  1 00:00:00 hostname sshd[000]: message"
    let actual = filterDate (line, today)
    Assert.Equal(true, actual)

[<Fact>]
let ``filterDate is not Ok`` () =
    let yesterday = today.AddDays(-1)
    let line = "Jan  1 00:00:00 hostname sshd[000]: message"
    let actual = filterDate (line, yesterday)
    Assert.Equal(false, actual)

[<Fact>]
let ``filterDate is not correct month`` () =
    let yesterday = today.AddDays(-1)
    let line = "Jan  1 00:00:00 hostname sshd[000]: message"
    let actual = filterDate (line, yesterday)
    Assert.Equal(false, actual)

[<Fact>]
let ``filterDate is not correct Day`` () =
    let yesterday =
        System.TimeZoneInfo.ConvertTime(new System.DateTime(2001, 1, 1), jstTimeZone)

    let line = "Jan  2 00:00:00 hostname sshd[000]: message"
    let actual = filterDate (line, yesterday)
    Assert.Equal(false, actual)

[<Fact>]
let ``filterDaemon daemon is correct (systemd-logind)`` () =
    let line = "systemd-logind[000]: message"
    let actual = filterDaemon (line)
    Assert.Equal(true, actual)

[<Fact>]
let ``filterDaemon daemon is crrect (sshd)`` () =
    let line = "sshd[000]: message"
    let actual = filterDaemon (line)
    Assert.Equal(true, actual)


[<Fact>]
let ``filterDaemon daemon is not correct (CRON)`` () =
    let line = "CRON[000]: message"
    let actual = filterDaemon (line)
    Assert.Equal(false, actual)

[<Fact>]
let ``filterLine daemon is correct (sshd)`` () =
    let line = "Jan  1 00:00:00 hostname sshd[000]: message"
    let actual = filterLine (line, 25, today)
    Assert.Equal(true, actual)

[<Fact>]
let ``filterLine daemon is correct (systemd-logind)`` () =
    let line = "Jan  1 00:00:00 hostname systemd-logind[000]: message"
    let actual = filterLine (line, 25, today)
    Assert.Equal(true, actual)

[<Fact>]
let ``filterLine daemon is not correct`` () =
    let line = "Jan  1 00:00:00 hostname CRON[000]: message"
    let actual = filterLine (line, 25, today)
    Assert.Equal(false, actual)

[<Fact>]
let ``filterLine month is not correct`` () =
    let yesterday = today.AddDays(-1)
    let line = "Jan  1 00:00:00 hostname sshd[000]: message"
    let actual = filterLine (line, 25, yesterday)
    Assert.Equal(false, actual)

[<Fact>]
let ``filterLine day is not correct`` () =
    let yesterday =
        System.TimeZoneInfo.ConvertTime(new System.DateTime(2001, 1, 1), jstTimeZone)

    let line = "Jan  2 00:00:00 hostname sshd[000]: message"
    let actual = filterLine (line, 25, yesterday)
    Assert.Equal(false, actual)
