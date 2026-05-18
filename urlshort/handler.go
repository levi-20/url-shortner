package urlshort

import (
	"log/slog"
	"net/http"

	"gopkg.in/yaml.v3"
)

type YamlRecords struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	parsedMYaml := []YamlRecords{}

	err := yaml.Unmarshal(yml, &parsedMYaml)
	if err != nil {
		slog.Error("error while parsing the yaml", "error", err)
		return nil, err
	}

	urlMappings := make(map[string]string)
	for _, mapping := range parsedMYaml {
		urlMappings[mapping.Path] = mapping.URL
	}

	return MapHandler(urlMappings, fallback), nil
}
