package gorbled

import (
    "log"

    "io/ioutil"
    "encoding/json"
)

const (
    CONFIG_FILE_PATH = "config.json"
)

type Config struct {
    Title, Description string
    PageSize int
}

/*
 * Get Config
 *
 * @return config (Config)
 */
func getConfig() (config Config) {
    configFile, err := ioutil.ReadFile(CONFIG_FILE_PATH)
    err = json.Unmarshal(configFile, &config)
    if err != nil {
        log.Fatal(err)
    }

    return
}
