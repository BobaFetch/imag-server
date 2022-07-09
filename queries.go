package main

import (
	"context"
	"fmt"
)

// department name & number
func getDepartmentInfo(ref string) Department {
	var dept Department

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish a connection: ", err.Error())
	}

	tsql := fmt.Sprintf(`
		SELECT WCNDESC, WCNREF 
		FROM WcntTable
		WHERE WCNREF = '%s'
	`, ref)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		fmt.Println("Error executing query: ", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&dept.Name,
			&dept.Number,
		)
		if err != nil {
			fmt.Println("Error getting department data: ", err.Error())
		}
	}
	return dept
}

// top 20 jobs in department
// sorted by priority, scheduled date
func getPartList(ref string) []Part {
	var job Part
	var jobList []Part

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish a conncetion: ", err.Error())
	}

	tsql := fmt.Sprintf(`
	SELECT DISTINCT TOP 20
		RUNREF,
		RUNRTNUM, 
		RUNNO,
		RUNQTY,
		SOCUST,
		--SOPO, 
		--RASOITEM, 
		ISNULL(AGPMCOMMENTS, '') AGPMCOMMENTS,
		RUNPRIORITY, 
		OPSCHEDDATE,
		ISNULL((SELECT DATEDIFF(MINUTE,(Select TOP 1 OPCOMPDATE From RnopTable WHERE OPREF = RUNREF AND OPRUN = RUNNO AND OPCOMPLETE IS NOT NULL ORDER BY OPCOMPDATE DESC),GETDATE())), 0) DTDIFF,
		WCNDESC
		FROM RunsTable
		INNER JOIN RnopTable ON OPREF=RUNREF AND OPRUN= RUNNO AND RUNOPCUR=OPNO
		INNER JOIN PartTable ON PARTREF=RUNREF 
		INNER JOIN RnalTable ON RUNREF=RAREF AND RUNNO=RARUN
		INNER JOIN SohdTable ON SONUMBER=RASO
		INNER JOIN WcntTable ON OPCENTER = WCNREF
		LEFT OUTER JOIN AgcmTable ON AGPART=RUNRTNUM AND AGRUN=RUNNO
		LEFT OUTER JOIN SoitTable ON ITPART=PARTREF AND ITSO=RASO

		WHERE OPCENTER LIKE '%s'
		AND OPCOMPLETE = 0
		ORDER BY RUNPRIORITY, OPSCHEDDATE`, ref)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		fmt.Println("Error executing query: ", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&job.PartRef,
			&job.PartNum,
			&job.Run,
			&job.Quantity,
			&job.Customer,
			&job.Comments,
			&job.Priority,
			&job.SchedDate,
			&job.QueueDiff,
			&job.WCName,
		)
		if err != nil {
			fmt.Println("Error getting jobs: ", err.Error())
		}
		jobList = append(jobList, job)
	}
	return jobList
}

// department daily goal
func getDailyGoal(ref string) int {
	var goal int

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish a connection: ", err.Error())
	}

	tsql := fmt.Sprintf(`
		SELECT 
				ISNULL((SELECT MAX(DTOTAL) daily_goal
			FROM (
				SELECT OPCENTER,
				ROW_NUMBER() OVER (PARTITION BY OPCENTER ORDER BY OPCENTER) DTOTAL 
				FROM RnopTable 
				WHERE OPCOMPLETE = 0 AND OPSCHEDDATE <= CAST(GETDATE() AS DATETIME) + 30 AND OPCENTER = '%s'
			)b GROUP BY OPCENTER), 0) daily_goal
	`, ref)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		fmt.Println("Error executing query: ", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&goal,
		)
		if err != nil {
			fmt.Println("Error getting data: ", err.Error())
		}
	}
	return goal
}

// completed jobs per department
func getCompletedJobCount(ref string) int {
	var completedJobs int

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish a connection: ", err.Error())
	}

	tsql := fmt.Sprintf(`
	SELECT 
		ISNULL((SELECT MAX(OPTOTAL)
	FROM (
		SELECT 
			OPCENTER, 
			ROW_NUMBER() OVER (PARTITION BY OPCENTER ORDER BY OPCENTER) OPTOTAL 
		FROM RnopTable 
		JOIN RunsTable ON OPREF=RUNREF AND OPRUN=RUNNO
		WHERE OPCOMPDATE >= CAST(GETDATE() AS DATE) AND OPCENTER = '%s' AND RUNSTATUS <> 'CA'
	)a 
	GROUP BY OPCENTER), 0) completed_jobs
	`, ref)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		fmt.Println("Error executing query: ", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&completedJobs,
		)
		if err != nil {
			fmt.Println("Error getting data: ", err.Error())
		}
	}
	return completedJobs
}

// completed parts per department
func getCompletedPartCount(ref string) int {
	var completedParts int

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish a connection: ", err.Error())
	}

	tsql := fmt.Sprintf(`
		SELECT 
			CAST(ISNULL(SUM(OPACCEPT), 0) AS INT) PART_COUNT
		FROM RnopTable
		WHERE OPCENTER LIKE '%s' AND OPCOMPDATE >= CAST(GETDATE() AS DATE)
	`, ref)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		fmt.Println("Error executing query: ", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&completedParts,
		)
		if err != nil {
			fmt.Println("Error getting data: ", err.Error())
		}
	}
	return completedParts
}

// employee daily statics per department
func getEmployeeDailyStats(ref string) []EmployeeStats {
	var dailyStats EmployeeStats
	var dailyStatsList []EmployeeStats

	ctx := context.Background()

	err := db.PingContext(ctx)
	if err != nil {
		fmt.Println("Could not establish a connection: ", err.Error())
	}

	tsql := fmt.Sprintf(`
		SELECT 
			OPINSP EMPLOYEE,
			MAX(OPTOTAL) COMPLETED_JOBS
		FROM(
			SELECT OPINSP,
			ROW_NUMBER() OVER (PARTITION BY OPINSP ORDER BY OPINSP) AS OPTOTAL
			FROM RnopTable
			JOIN RunsTable ON OPREF = RUNREF AND OPRUN = RUNNO
			WHERE OPCOMPDATE >= CAST(GETDATE() AS DATE)
			AND RUNSTATUS <> 'CA'
			AND OPCENTER LIKE '%s')a
		GROUP BY OPINSP
	`, ref)

	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		fmt.Println("Error executing query: ", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&dailyStats.Employee,
			&dailyStats.CompletedJobs,
		)
		if err != nil {
			fmt.Println("Error getting data: ", err.Error())
		}
		dailyStatsList = append(dailyStatsList, dailyStats)
	}
	return dailyStatsList
}
