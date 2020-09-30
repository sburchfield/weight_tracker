package main

var queries = map[string]string{

	"getTodayCalories": `SELECT * FROM public.calorie_tracker where calorie_tracker.updated_at like to_timestamp(concat(current_date , ' 08:00:00'), 'YYYY-MM-DD hh24:mi:ss')`,

	"getOverallCalories": `SELECT SUM(calories) AS total, DATE(created_at) as date FROM calories WHERE user_id = ? GROUP BY DATE(created_at);`,
}
