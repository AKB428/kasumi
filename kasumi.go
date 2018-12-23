package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

// Conf ... ConoHa API sアクセス情報
type Conf struct {
	AuthURL    string `json:"auth_url"`
	TenantName string `json:"tenantName"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	EndPoint   string `json:"endPoint"`
}

// AuthToken ... AuthT APIのレスポンスJSONを定義
type AuthToken struct {
	Access struct {
		Token struct {
			IssuedAt string    `json:"issued_at"`
			Expires  time.Time `json:"expires"`
			ID       string    `json:"id"`
			Tenant   struct {
				Description string `json:"description"`
				Enabled     bool   `json:"enabled"`
				ID          string `json:"id"`
				Name        string `json:"name"`
			} `json:"tenant"`
		} `json:"token"`
		User struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
		ServiceCatalog []struct {
			Endpoints []struct {
				ID          string `json:"id"`
				AdminURL    string `json:"adminURL"`
				InternalURL string `json:"internalURL"`
				PublicURL   string `json:"publicURL"`
				Region      string `json:"region"`
			} `json:"endpoints"`
			EndpointsLinks []interface{} `json:"endpoints_links"`
			Type           string        `json:"type"`
			Name           string        `json:"name"`
		} `json:"serviceCatalog"`
	} `json:"access"`
}

var wg sync.WaitGroup
var glNum = 200

func main() {

	flag.Parse()
	containerName := flag.Arg(0)
	fmt.Print(containerName)

	const format = "20060102_150405"
	logFileName := "./log/" + time.Now().Format(format) + ".log"

	logFile, _ := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY, 0666)

	defer logFile.Close()
	log.SetOutput(logFile)

	bytes, err := ioutil.ReadFile("./conf/conoha_api_v1_key.json")
	if err != nil {
		log.Fatal(err)
	}
	// JSONデコード
	var conf Conf

	if err := json.Unmarshal(bytes, &conf); err != nil {
		log.Fatal(err)
	}
	// デコードしたデータを表示
	fmt.Printf("%+v\n", conf)

	token := getToken(conf)

	fmt.Println(token)

	deleteFileCounter := 0

	for {
		// 指定されたフォルダを再帰的にリスト取得 <リストAPI>
		objectList := getContainerList(token, conf.EndPoint, containerName)

		log.Println(fmt.Sprintf("%s: %d", "objectList size", len(objectList)))

		if len(objectList) == 0 {
			break
		}

		// goroutineで100ぐらい一気に削除　　<削除API>
		counter := 0

		for _, url := range objectList {

			wg.Add(1)
			go deleteObject(token, url)

			deleteFileCounter++
			counter++
			if counter == glNum {
				counter = 0
				wg.Wait()
				log.Println("wait done")
				log.Println(fmt.Sprintf("%s: %d", "deleteFileCounter", deleteFileCounter))
			}
		}

		// エラーは無視して順次削除を繰り返す
	}

	// TODO 処理時間カウント
	// TODO 処理ファイル数カウント
}

//# curl -i 'https://********.jp/v2.0/tokens' -X POST -H "Content-Type: application/json" -H "Accept: application/json"  -d '{"auth": {"tenantName": "1234567", "passwordCredentials": {"username": "1234567", "password": "************"}}}'
func getToken(conf Conf) string {
	//jsonStr := `{"tenantName":"` + tenatName + `","device":"` + device + `"}`

	jsonStr := `{"auth": {"tenantName": "` + conf.TenantName + `", "passwordCredentials": {"username": "` + conf.Username + `", "password": "` + conf.Password + `"}}}`

	req, err := http.NewRequest(
		"POST",
		conf.AuthURL,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		fmt.Println(err)
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	fmt.Println(response.Status)

	// body, error
	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(body))

	//jsonBytes := ([]byte)(string(body))
	data := new(AuthToken)

	if err := json.Unmarshal(body, data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
	}

	fmt.Printf("%+v\n", data)

	fmt.Println(data.Access.Token.ID)

	token := data.Access.Token.ID

	//TODO swiftのURLはレスポンスから取得する
	return token
}

func getContainerList(token string, baseURL string, containerName string) []string {

	url := baseURL + containerName

	fmt.Println(url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Auth-Token", token)

	client := new(http.Client)
	response, err := client.Do(req)

	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	body, _ := ioutil.ReadAll(response.Body)

	fmt.Print(string(body))

	var objectList []string
	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(string(body), -1) {

		if v != "" {
			//fmt.Println(i+1, ":", url + "/" +v)

			objectList = append(objectList, url+"/"+v)
		}
	}
	return objectList
}

func deleteObject(token string, url string) {
	defer wg.Done()

	// fmt.Println("DELETE: " + url)

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("X-Auth-Token", token)

	client := new(http.Client)
	response, err := client.Do(req)

	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	//body, _ := ioutil.ReadAll(response.Body)

	//fmt.Println(response.Status) string 204 No Content

	if response.StatusCode != 204 {
		fmt.Println(response.Status)
	}

}
