## 로그 분석기


### S3 에 저장되는 하루치 로그를 간단하게 분석하기

로그 파일은 하루치만 보았을때도 정말 많은 로그들이 쌓이기 마련입니다.


그래서 큰 용량과 많은 파일을 다루기위해 Golang 의 고루틴과 채널을 사용하여 처리하는 로직을 개발하였습니다.


### chan(채널) 을 사용한 이유
1. 동시성 처리
    - 대량의 로그 파일을 병렬로 처리하면서 성능을 극대화합니다. 각 파일에 대한 파싱 작업이 고루틴을 통해 병렬로 실행되므로 파일이 많더라도 짧은 시간 안에 데이터를 수집할 수 있습니다.
   
2. 안전한 데이터 공유
   - Go에서는 여러 고루틴이 동시에 데이터를 다룰 때 경쟁 상태(race condition)가 발생할 수 있습니다. 채널은 이러한 상태를 피하면서 데이터를 주고받을 수 있는 안전한 방법을 제공합니다. 따라서, 고루틴 간 데이터가 안전하게 전달되므로, 예기치 않은 오류를 방지할 수 있습니다.
3. 작업 흐름 관리
   - 채널을 통해 데이터를 전달함으로써, 모든 파일이 처리될 때까지 메인 루틴이 기다릴 수 있도록 합니다. 이때, sync.WaitGroup을 함께 사용하여 모든 고루틴이 작업을 완료할 때까지 대기할 수 있습니다. 이로써, 병렬 작업이 완료된 후에만 분석 결과를 출력하게 됩니다.


### 파일 구성
- analysis/log-analysis.go - 로그 데이터를 처리하여 요약
- parser/log-parser.go - 읽은 파일의 각 로그별로 LogEntry 모델로 변환
- processors/file-processor.go - .gz확장자 로그 파일을 처리하여 각 항목을 파싱
- utils/read-zip_file.go - 병렬로 여러개의 gzip 파일을 읽고 처리



### 주요 로직
```go

func MultipleReadGzipFile(dirPath string) ([]models.LogEntry, error) {
    // 나머지 로직
    
    resultChan := make(chan models.LogEntry, 100)
    var wg sync.WaitGroup
    
    // 나머지 로직
}


func LogAnalysis(entries []models.LogEntry) models.LogAnalysis {
    analysis := models.LogAnalysis{
        LevelCounts:     make(map[string]int),
        MethodCounts:    make(map[string]int),
        MethodDurations: make(map[string]float64),
    }

    // 나머지 로직
}


```
make를 사용해 채널을 생성하면, 고루틴 간에 데이터를 전달하는데 필요한 채널을 지정된 버퍼 크기로 설정할 수 있습니다.
해당 설정을 통해 비동기적으로 데이터를 주고받을 수 있어 성능 향상에 도움이 됩니다.

- resultChan 에서는 버퍼 크기를 100으로 초기화 해, 고루틴간에 발생하는 로그 항목을 최대 100까지 대기열에 쌓아두고 처리할 수 있도록 했습니다.
- map 초기화에서 make을 사용했습니다. 이는 값 할당 오류를 방지하고, 메모리를 효율적으로 사용할 수 있게 합니다.


