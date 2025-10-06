package usecases

import (
	"bytes"
	user "eaglechat/apps/id_manager/internal/domain/repositories/user"
	"encoding/json"
	"log"
	"net"
	"net/http"
)

// SyncUsersUseCase handles peer synchronization.
type SyncUsersUseCase struct {
	repo user.UserRepository
}

// NewSyncUsersUseCase creates a new instance of SyncUsersUseCase.
func NewSyncUsersUseCase(repo user.UserRepository) *SyncUsersUseCase {
	return &SyncUsersUseCase{repo: repo}
}

// SyncWithPeer requests user data from a peer and merges it.
func (uc *SyncUsersUseCase) SyncWithPeer(peerIP string) {
	localUsers, err := uc.repo.FindAll()
	if err != nil {
		log.Println("Error getting local users:", err)
		return
	}

	// Get local ids to send to peer
	usernames := make([]string, 0, len(localUsers))
	for _, u := range localUsers {
		usernames = append(usernames, u.Username)
	}

	reqBody := map[string]interface{}{
		"usernames":         usernames,
		"omit_disconnected": false,
	}

	payload, _ := json.Marshal(reqBody)

	// Make a request to the peer
	resp, err := http.Post("http://"+peerIP+":8080/users/query", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Error connecting to peer:", err)
		return
	}
	defer resp.Body.Close()

	var peerData map[string]*QueryUserResponseData
	if err := json.NewDecoder(resp.Body).Decode(&peerData); err != nil {
		log.Println("Error decoding peer response:", err)
		return
	}

	// Update data in current server
	for Id, data := range peerData {
		localUser, err := uc.repo.FindByID(Id)
		if err != nil {
			// if user does not exist locally we do not update it
			// we ignore it for now
			// TODO: make endpoint to get every user and update with
			// current id manager
			continue
		}

		// update mutable fields
		localUser.PublicKeyPEM = data.PublicKey
		localUser.IP = data.IP

		if err := uc.repo.Save(localUser); err != nil {
			log.Printf("Error updating user %s: %v", Id, err)
		}
	}
}

type QueryUserResponseData struct {
	Id        string  `json:"Id"`
	PublicKey string  `json:"public_key"`
	IP        *net.IP `json:"ip,omitempty"`
}
