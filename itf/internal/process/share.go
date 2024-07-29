package process

import (
	"fmt"
	"log"
	"os"

	"github.com/apaichon/thesis-itf/itf/internal/models"
	"github.com/apaichon/thesis-itf/itf/internal/repositories"
	"github.com/google/uuid"
)

type Status struct {
	MessageID uuid.UUID `json:"messageId"`
	Error     string    `json:"error"`
}


func TransferIntraBankRestApiFlow(task models.TaskQueueModel) {
	statusList := TransferIntraBankRestFlow(task)
	WriteStatusListToFile(statusList, "status.log")
}

func TransferIntraBankFlow(task models.TaskQueueModel) error {

	repo := repositories.NewProcessManagerRepo()
	repo.FlagDeleteTaskQueue(task.Id)
	messages, err := repo.GetMessagesByIds(task.MessageIds)
	if err != nil {
		fmt.Printf("Error:%v", err)
		return err
	}

	log.Printf("len Transfer Message:%v", len(messages))

	banking360 := NewBanking360()
	query, variables := banking360.TransformIntraBank(messages)
	banking360.PostToGraphQL(query, variables)

	log.Println("Transfer Intra Bank Successfully!")
	return nil
}

func TransferIntraBankRestFlow(task models.TaskQueueModel) []Status {

	var statusList []Status
	repo := repositories.NewProcessManagerRepo()
	repo.FlagDeleteTaskQueue(task.Id)
	messages, err := repo.GetMessagesByIds(task.MessageIds)
	if err != nil {
		log.Printf("Get Message Ids: %v", err)
		statusList = append(statusList, Status{Error: err.Error()})
		return statusList
	}

	log.Printf("len Transfer Message:%v", len(messages))

	banking360 := NewBanking360()
	for _, message := range messages {
		transfer, err := banking360.ConvertMessageToTransferInput(message)
		if err != nil {
			log.Println("Convert Transfer Intra Error!")
			statusList = append(statusList, Status{MessageID: message.Id, Error: err.Error()})
			continue
		}

		err = banking360.PostTransferIntraBankInput(transfer)
		if err != nil {
			log.Println("Post Transfer Intra Error!")
			statusList = append(statusList, Status{MessageID: message.Id, Error: err.Error()})
			continue
		}

		log.Println("Transfer Intra Bank Successfully!")

	}
	return statusList
}

func WithdrawalCoreBankFlow(task models.TaskQueueModel) error {

	repo := repositories.NewProcessManagerRepo()
	repo.FlagDeleteTaskQueue(task.Id)
	messages, err := repo.GetMessagesByIds(task.MessageIds)
	if err != nil {
		fmt.Printf("Error:%v", err)
		return err
	}

	log.Printf("len withdraw Message:%v", len(messages))

	corebanking := NewCoreBanking()
	for _, message := range messages {
		withdrawal, err := corebanking.TransformWithdrawal(message)

		if err == nil {

			result, err := corebanking.PostWithdrawal(*withdrawal)
			if err != nil {
				fmt.Errorf("[error]Withdrawal:%v", err)
			}
			log.Println("Withdrawal result:" + result)

		}
	}

	log.Println("Withdrawal Core Bank Successfully!")
	return nil
}

func DepositCoreBankFlow(task models.TaskQueueModel) error {

	repo := repositories.NewProcessManagerRepo()
	repo.FlagDeleteTaskQueue(task.Id)
	messages, err := repo.GetMessagesByIds(task.MessageIds)
	if err != nil {
		fmt.Printf("Error:%v", err)
		return err
	}

	log.Printf("len transfer Message:%v", len(messages))

	corebanking := NewCoreBanking()
	for _, message := range messages {

		withdrawal, err := corebanking.TransformDeposit(message)

		if err == nil {

			result, err := corebanking.PostDeposit(*withdrawal)
			if err != nil {
				fmt.Errorf("[error]Deposit:%v", err)
			}
			log.Println("Deposit result:" + result)

		}
	}

	log.Println("Deposit Core Bank Successfully!")
	return nil
}

func WriteStatusListToFile(statusList []Status, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, status := range statusList {
		statusLog := fmt.Sprintf("Message ID: %s, Error: %s\n", status.MessageID, status.Error)
		if _, err := file.WriteString(statusLog); err != nil {
			return err
		}
	}

	return nil
}
