package quiz

func canStart(quiz *Quiz) bool {
	if len(quiz.Title) == 0 {
		return false
	}
	if len(quiz.Questions) == 0 {
		return false
	}
	for _, q := range quiz.Questions {
		if !isQuestionValid(&q) {
			return false
		}
	}
	return true
}
func isQuestionValid(question *Question) bool {
	if len(question.Title) == 0 {
		return false
	}
	if len(question.Answers) < 2 {
		return false
	}
	corrects := 0
	for _, answer := range question.Answers {
		if answer.IsCorrect {
			corrects++
		}
	}
	if corrects == len(question.Answers) {
		return false
	}
	return true
}
