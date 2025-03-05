package quizzes

type dummyEntry struct {
	ownerId string
	quizzes []Quiz
}

type DummyQuizStoreImpl struct {
	entries []dummyEntry
}

func (d DummyQuizStoreImpl) Upsert(ownerId string, quiz Quiz) error {
	for _, u := range d.entries {
		if u.ownerId == ownerId {
			for _, q := range u.quizzes {
				if q.Id == quiz.Id {
					q.Title = quiz.Title
					q.Description = quiz.Description
					q.Questions = quiz.Questions
					q.Code = quiz.Code
					return nil
				}
			}

			u.quizzes = append(u.quizzes, quiz)
		}
	}

	return nil
}

func (d DummyQuizStoreImpl) GetUnique(ownerId, uid string) (Quiz, error) {
	for _, u := range d.entries {
		if u.ownerId == ownerId {
			for _, q := range u.quizzes {
				if q.Id == uid {
					return q, nil
				}
			}
		}
	}

	return Quiz{}, ErrNotFound
}

func (d DummyQuizStoreImpl) GetQuizzes(ownerId string) ([]Quiz, error) {
	for _, u := range d.entries {
		if u.ownerId == ownerId {
			return u.quizzes, nil
		}
	}

	return []Quiz{}, ErrNotFound
}

func (d DummyQuizStoreImpl) Patch(ownerId, uid string, fields []FieldPatchOp) error {
	quiz, err := d.GetUnique(ownerId, uid)
	if err != nil {
		return err
	}

	for _, field := range fields {
		if field.Op == "replace" {
			if field.Path == "/title" {
				quiz.Title = field.Value.(string)
			} else {
				return ErrInvalidPatchField
			}
		} else {
			return ErrInvalidPatchOperator
		}
	}

	return nil
}

func (d DummyQuizStoreImpl) GetUniqueQuestion(ownerId, quizId, questionId string) (Question, error) {
	for _, u := range d.entries {
		if u.ownerId == ownerId {
			for _, q := range u.quizzes {
				if q.Id == quizId {
					for _, qu := range q.Questions {
						if qu.Id == questionId {
							return qu, nil
						}
					}
				}
			}
		}
	}

	return Question{}, ErrNotFound
}

func (d DummyQuizStoreImpl) UpsertQuestion(ownerId, quizId string, question Question) error {
	for _, u := range d.entries {
		if u.ownerId == ownerId {
			for _, q := range u.quizzes {
				if q.Id == quizId {
					for _, qu := range q.Questions {
						if qu.Id == question.Id {
							qu.Title = question.Title
							qu.Answers = question.Answers
							return nil
						}
					}

					q.Questions = append(q.Questions, question)
					return nil
				}
			}
		}
	}

	return nil
}

func (d DummyQuizStoreImpl) UpdateQuestion(ownerId, quizId string, question Question) error {
	for _, u := range d.entries {
		if u.ownerId == ownerId {
			for _, q := range u.quizzes {
				if q.Id == quizId {
					for _, qu := range q.Questions {
						if qu.Id == question.Id {
							qu.Title = question.Title
							qu.Answers = question.Answers
							return nil
						}
					}

					return ErrNotFound
				}
			}
		}
	}

	return ErrNotFound
}

func CreateDummyStore() Store {
	return DummyQuizStoreImpl{
		entries: make([]dummyEntry, 0),
	}
}
