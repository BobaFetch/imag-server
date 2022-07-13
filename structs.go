package main

type Part struct {
	PartRef   string  `json:"partRef"`
	PartNum   string  `json:"partNumber"`
	Run       string  `json:"run"`
	Quantity  float32 `json:"qty"`
	Customer  *string `json:"customer"`
	Comments  *string `json:"comments"`
	Priority  int     `json:"priority"`
	SchedDate *string `json:"schedDate"`
	QueueDiff int     `json:"queueDiff"`
	WCName    string  `json:"wcName"`
}

type Department struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

type QueueData struct {
	Department Department      `json:"department"`
	Parts      []Part          `json:"parts"`
	Stats      DepartmentStats `json:"stats"`
	Employees  []EmployeeStats `json:"employeeStats"`
}

type DepartmentStats struct {
	Goal           int           `json:"dailyGoal"`
	CompletedJobs  int           `json:"completedJobs"`
	CompletedParts int           `json:"completedParts"`
	WeeklyStats    []WeeklyStats `json:"weeklyStats"`
}

type EmployeeStats struct {
	Employee      string `json:"employee"`
	CompletedJobs int    `json:"completedJobs"`
}

type WeeklyStats struct {
	CompletedJobs  int    `json:"completedJobs"`
	CompletionDate string `json:"completionDate"`
}

type Burndown struct {
	Part_Num   string  `json:"partNumber"`
	Run        string  `json:"run"`
	Comments   string  `json:"comments"`
	WC_Num     string  `json:"wcNumber"`
	WC_Name    string  `json:"wcName"`
	Quantity   float32 `json:"quantity"`
	Queue_Diff int     `json:"queueDiff"`
}
