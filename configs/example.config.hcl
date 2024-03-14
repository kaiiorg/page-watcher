page {
  name = "Kaiiorg home page"
  url = "https://kaiiorg.wtf/"
  every = "30s"

  find = ["main", "class", "main-wrap"]

  normalize = {
    // Consecutive spaces to single space
    "[[:blank:]]+" = " "
    // Consecutive NBSP to single space
    " +" = " "
    // Blank lines to empty line
    "[[:blank:]]+\n" = ""
  }

  debug = false
}
