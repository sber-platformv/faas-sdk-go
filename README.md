# FAAS-SDK для тестирования и локальной отладки Golang функций

FAAS-SDK позволяет локально тестировать функции [Platform V Functions](https://developers.sber.ru/portal/products/platform-v-functions) без необходимости писать HTTP сервер и логику обработки запросов.

SDK для Golang устанавливается как обычный пакет.

## Пререквизиты

- Golang 1.17.x;

## Установка и использование

1. Перейдите в директорию проекта с функцией. Если вы только начинаете локальную разработку, вы можете:
   * Экспортировать функцию из Functions и распаковать архив в произвольную локальную директорию.
   * Начать работу с чистого проекта - для этого создайте файл `handlers/handler.go` содержаший:

   		```golang
    	 package handlers
		
    	 import (
        	 "fmt"  
        	 "io/ioutil"
        	 "log"
        	 "net/http"
    	 )
		 
    	 // Метод Handler. Данный метод будет обрабатывать HTTP запросы поступающие к функции
    	 func Handler(w http.ResponseWriter, r *http.Request) {
        	 if r.Body != nil {
            	 defer r.Body.Close()
        	 }
			 
        	 // Чтение тела запроса
        	 body, _ := ioutil.ReadAll(r.Body)
			 
        	 // Логирование входящего запроса
        	 log.Printf("Request received: %s\nMethod: %s", string(body), r.Method)
			 
        	 // Подготовка и возврат ответа на вызов
        	 w.WriteHeader(http.StatusOK)
        	 w.Write([]byte(fmt.Sprintf("Hello from Go function!\nYou said: %s", string(body))))
    		 }
    	```

2. Если ваш проект с функцией еще не инициирован - создайте модуль Go. Для этого запустите CLI и перейдите в нем в директорию проекта с функцией. Затем выполните команду:
    
    ```shell
     go mod init example.com/test
    ```
    > Вместо `example.com/test` можно использовать любое другое имя для модуля.

3. Создайте либо откройте файл `main.go`, и добавьте в него следующие строки:
    
    ```go
    package main

    import (
        "example.com/test/handlers" // Ваш пакет с функцией

        "github.com/sber-platformv/faas-sdk-go/framework" // Пакет с SDK
    )

    func main() {
        port := "8080" // Порт, на котором будет поднят слушатель сервера для получения запросов к функции
        framework.Start(port, handlers.Handler) // метод принимает в аргументы порт и метод функции, которую нужно запустить
    }
    ```

4. Обновите зависимости проекта через CLI командой:
    
    ```shell
    go mod tidy
    ```

5. Запустите сервер командой в CLI:
    
    ```shell
    go run main.go
    ```

6. Отправляйте запросы используя `curl`, браузер или другие инструменты:
    
    ```shell
    curl localhost:8080
    # Hello from Go function!
    # You said:
    ```

## Конфигурация

Запуск SDK конфигурируется с помощью настроек в файле `main.go`. Вы можете изменить следующие параметры:
* `port` -  Порт, на котором будет поднят слушатель сервера для получения запросов к функции.

## Unit-тестирование

Вы можете добавить unit-тесты в тестируемую локально функцию так же, как и в любой Golang проект.

Например, для тестирования `hello, world` примера:

1. Создайте в директории `handlers` файл `handler_test.go` со следующим содержанием:
    
    ```go
    package handlers

    import (
        "io/ioutil"
        "net/http/httptest"
        "testing"
    )

    func Test_Handler(t *testing.T) {
        w := httptest.NewRecorder()
        r := httptest.NewRequest("GET", "/", nil)

        Handler(w, r)
        body, _ := ioutil.ReadAll(w.Body)
        expected := "Hello from Go function!\nYou said: "
        if string(body) != expected {
            t.Error(
                "expected", expected,
                "got", string(body),
            )
        }
    }
    ```
2. Перейдите в директорию проекта в CLI и запустите тесты командой `go test ./...`. Результаты теста будут отображены в терминале.

