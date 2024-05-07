module Tests

open Xunit

open Program

let jstTimeZone = System.TimeZoneInfo.FindSystemTimeZoneById("Tokyo Standard Time")
let format = System.Globalization.CultureInfo.CreateSpecificCulture("en-US")
let today = System.TimeZoneInfo.ConvertTime(System.DateTime.Today, jstTimeZone)

[<Fact>]
let ``daemon is sshd && today`` () =
    let line =
        today.ToString("MMM", format)
        + today.Day.ToString().PadLeft(3)
        + " 00:00:00 hostname sshd[000]: message"

    let actual = filterLine (line, 25)
    Assert.Equal(true, actual)

[<Fact>]
let ``daemon is systemd-logind && today`` () =
    let line =
        today.ToString("MMM", format)
        + today.Day.ToString().PadLeft(3)
        + " 00:00:00 hostname systemd-logind[000]: message"

    let actual = filterLine (line, 25)
    Assert.Equal(true, actual)

[<Fact>]
let ``daemon is not target`` () =
    let line =
        today.ToString("MMM", format)
        + today.Day.ToString().PadLeft(3)
        + " 00:00:00 hostname CRON[000]: message"

    let actual = filterLine (line, 25)
    Assert.Equal(false, actual)

[<Fact>]
let ``date is not target`` () =
    let yesterday = today.AddDays(-1)

    let line =
        yesterday.ToString("MMM", format)
        + yesterday.Day.ToString().PadLeft(3)
        + " 00:00:00 hostname sshd[000]: message"

    let actual = filterLine (line, 25)
    Assert.Equal(false, actual)
