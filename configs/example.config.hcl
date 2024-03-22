page {
  name  = "page-watcher demo page"
  url   = "http://localhost:8080/demo"
  every = "30s"

  find = ["p", "id", "target"]

  normalize {
    // Consecutive NBSP to single space
    regex = " +"
    to    = " "
  }
  normalize {
    // Consecutive spaces to single space
    regex = "[[:blank:]]+"
    to    = " "
  }
  normalize {
    // Blank lines to empty line
    regex = "[[:blank:]]+\n"
    to    = ""
  }

  debug = false
}

db {
  path   = "./page-watcher.db"
  retain = 5
}

web {
  port = 8080
}
