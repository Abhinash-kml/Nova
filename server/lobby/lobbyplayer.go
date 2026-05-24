package lobby

func (p *LobbyPlayer) SetState(new LobbyPlayerState) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.State = new
}

func (p *LobbyPlayer) SetName(newName string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.UserName = newName
}
