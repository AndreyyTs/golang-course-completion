package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func MyExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	var chIn, chOut chan interface{}
	for _, val := range jobs {
		chIn = chOut
		chOut = make(chan interface{}, 20)
		wg.Add(1)
		go func(j job, chIn, chOut chan interface{}) {
			j(chIn, chOut)
			close(chOut) //закрываем выходной канал после работы job'ы
			wg.Done()
		}(val, chIn, chOut)
	}
	wg.Wait() // ожидаем завершения горутин
}

func MySingleHash(in, out chan interface{}) {
	mytexDataSignerMd5 := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	for i := range in {
		wg.Add(1)
		data, _ := i.(int)
		go func(data int) {
			str := strconv.Itoa(data)
			md5 := make(chan string, 1)
			go func(ch chan string, str string) {
				mytexDataSignerMd5.Lock()
				ch <- DataSignerMd5(str)
				mytexDataSignerMd5.Unlock()

			}(md5, str)

			crc32_1 := make(chan string, 1)
			go func(ch chan string, str string) {
				ch <- DataSignerCrc32(str)
			}(crc32_1, str)

			crc32_md5 := make(chan string, 1)
			go func(ch, in chan string) {
				ch <- DataSignerCrc32(<-in)
			}(crc32_md5, md5)

			go func(ch1, ch2 chan string, o chan interface{}) {
				crc32 := <-ch1
				crc32_md5 := <-ch2
				o <- crc32 + "~" + crc32_md5
				wg.Done()
			}(crc32_1, crc32_md5, out)

		}(data)
		//runtime.Gosched()
	}
	wg.Wait()
}

func MyMultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for data := range in {
		wg.Add(1)
		go func(out chan interface{}, data interface{}) {
			resultsMuxtex := &sync.Mutex{}
			results := make([]string, 6, 6)
			wgRes := &sync.WaitGroup{}
			for i := 0; i <= 5; i++ {
				wgRes.Add(1)
				go func(idx int) {
					str := DataSignerCrc32(strconv.Itoa(idx) + data.(string))
					resultsMuxtex.Lock()
					results[idx] = str
					resultsMuxtex.Unlock()
					wgRes.Done()
				}(i)
			}
			wgRes.Wait()
			sssss := strings.Join(results, "")
			out <- sssss
			wg.Done()
		}(out, data)

		//runtime.Gosched()
	}
	wg.Wait()
}

func MyCombineResults(in, out chan interface{}) {
	results := make([]string, 0, 20)
	for i := range in {
		results = append(results, i.(string))
	}

	sort.Strings(results)
	out <- strings.Join(results, "_")
}

func main() {

	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	//inputData := []int{0, 1}
	//testResult := "NOT_SET"

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {

				out <- fibNum
				//runtime.Gosched()
			}
		}),
		job(MySingleHash),
		job(MyMultiHash),
		job(MyCombineResults),
		job(func(in, out chan interface{}) {

			for i := range in {
				strr, _ := i.(string)
				fmt.Println("testResult: ", strr)
				//runtime.Gosched()
			}

		}),
	}

	MyExecutePipeline(hashSignJobs...)
	//fmt.Scanln()

}
