page {
  name = "Kaiiorg home page"
  url = "https://kaiiorg.wtf/"
  every = "30s"

  find = ["main", "class", "main-wrap"]

  normalize {
    // Consecutive spaces to single space
    regex = " +"
    to = " "
  }
  normalize {
    // Consecutive NBSP to single space
    regex = "[[:blank:]]+"
    to = " "
  }
  normalize {
    // Blank lines to empty line
    regex = "[[:blank:]]+\n"
    to = ""
  }

  debug = false
}
