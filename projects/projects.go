package projects

import "strings"

type ProjectStats struct {
    accuracy float32
    precision float32
}

type Project struct {
    Name string `json:"name"`
    Model string `json:"model"`
    DataLocation string `json:"dataLocation"`
    URL string `json:"url"`
    Stats ProjectStats `json:"statistics"`
}

func getProjectKey(name string) string { return "project:" + strings.ToLower(name) }
func (P *Project) GetKey() string {
    return getProjectKey(P.Name)
}

func (p *Project) SetFromData(data map[string]string) {
    p.Name = data["name"]
    p.Model = data["model"]
    p.DataLocation = data["dataLocation"]
    p.URL = data["url"]
}

func (p *Project) SetStats(s ProjectStats) {
    p.Stats = s
}

func (p *Project) ToMap() map[string]string {
    var m map[string]string

    m["name"] = p.Name
    m["model"] = p.Model
    m["dataLocation"] = p.DataLocation
    m["url"] = p.URL

    return m
}
