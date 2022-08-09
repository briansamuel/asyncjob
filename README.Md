# AsyncJob
> move by briansamuel/asyncjob helper use concurrent to managerment job test.

[![Go Reference](https://pkg.go.dev/badge/github.com/princjef/gomarkdoc.svg)](https://pkg.go.dev/github.com/briansamuel/asyncjob)

### Install

``` bash
 go get github.com/briansamuel/asyncjob
```

### Usage

Library have 2 modules job and groupJob.
For example to setup:

* Add New Job

``` go
job1 := asyncjob.NewJob(func(ctx context.Context) error {
	time.Sleep(time.Second)
	log.Println("I am job 1")

	return nil
	
})
```


* New Job With Name


``` go

job3 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second * 2)
		log.Println("I am job 3")

		return nil
	}, asyncjob.WithName("Somthing with Job 3"))


```


* Add job to group
``` go
	group := asyncjob.NewGroup(true, job1, job2, job3)

	if err := group.Run(context.Background()); err != nil {
		log.Println(err)
	}
```
