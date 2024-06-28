package assets

import (
	"bytes"
	_ "embed"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	defaultSFXVolume   = 1.0
	defaultMusicVolume = 1.0

	sfxSampleRate = 44100
	//musicSampleRate = 192000
)

var (
	ctx   = audio.NewContext(sfxSampleRate)
	music *audio.Player

	// Music

	//go:embed audio/boss0_music.ogg
	boss0MusicBytes  []byte
	boss0MusicPlayer *audio.Player
	//go:embed audio/menu_music.ogg
	menuMusicBytes  []byte
	menuMusicPlayer *audio.Player

	// SFX

	//go:embed audio/shoot.wav
	shootSoundBytes  []byte
	shootSoundPlayer *audio.Player
	//go:embed audio/shoot2.wav
	shoot2SoundBytes  []byte
	shoot2SoundPlayer *audio.Player
	//go:embed audio/shoot3.wav
	shoot3SoundBytes  []byte
	shoot3SoundPlayer *audio.Player
	shootPlayers      [3]*audio.Player

	//go:embed audio/miss.wav
	missSoundBytes  []byte
	missSoundPlayer *audio.Player
	//go:embed audio/dash.wav
	dashSoundBytes  []byte
	dashSoundPlayer *audio.Player
	//go:embed audio/hit.wav
	hitSoundBytes  []byte
	hitSoundPlayer *audio.Player

	//go:embed audio/sm_shoot.wav
	smShootSoundBytes  []byte
	smShootSoundPlayer *audio.Player
	//go:embed audio/sm_comet.wav
	smCometSoundBytes  []byte
	smCometSoundPlayer *audio.Player
	//go:embed audio/boss_charge.wav
	bossChargeSoundBytes  []byte
	bossChargeSoundPlayer *audio.Player

	//go:embed audio/bonus.wav
	bonusSoundBytes  []byte
	bonusSoundPlayer *audio.Player
	//go:embed audio/portal.wav
	portalSoundBytes  []byte
	portalSoundPlayer *audio.Player
)

func newSFXPlayer(data []byte) *audio.Player {
	wavReader, err := wav.DecodeWithSampleRate(sfxSampleRate, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	p, err := ctx.NewPlayer(wavReader)
	if err != nil {
		log.Fatal(err)
	}

	return p
}

func newMusicPlayer(data []byte) *audio.Player {
	oggReader, err := vorbis.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	p, err := ctx.NewPlayer(audio.NewInfiniteLoop(oggReader, oggReader.Length()))
	if err != nil {
		log.Fatal(err)
	}

	return p
}

func init() {
	shootSoundPlayer = newSFXPlayer(shootSoundBytes)
	shoot2SoundPlayer = newSFXPlayer(shoot2SoundBytes)
	shoot3SoundPlayer = newSFXPlayer(shoot3SoundBytes)
	shootPlayers = [3]*audio.Player{
		shootSoundPlayer,
		shoot2SoundPlayer,
		shoot3SoundPlayer,
	}
	for i := range shootPlayers {
		shootPlayers[i].SetVolume(0.4)
	}

	missSoundPlayer = newSFXPlayer(missSoundBytes)
	dashSoundPlayer = newSFXPlayer(dashSoundBytes)
	dashSoundPlayer.SetVolume(0.5)
	hitSoundPlayer = newSFXPlayer(hitSoundBytes)
	//hitSoundPlayer.SetVolume(0.5)

	smShootSoundPlayer = newSFXPlayer(smShootSoundBytes)
	smCometSoundPlayer = newSFXPlayer(smCometSoundBytes)
	bossChargeSoundPlayer = newSFXPlayer(bossChargeSoundBytes)
	bonusSoundPlayer = newSFXPlayer(bonusSoundBytes)
	bonusSoundPlayer.SetVolume(0.5)
	portalSoundPlayer = newSFXPlayer(portalSoundBytes)
	portalSoundPlayer.SetVolume(0.5)

	boss0MusicPlayer = newMusicPlayer(boss0MusicBytes)
	menuMusicPlayer = newMusicPlayer(menuMusicBytes)
}

// Musics

/*func ReplayGameMusic() {
	gameMusicPlayer.Rewind()
	gameMusicPlayer.Play()
}

func StopGameMusic() {
	if gameMusicPlayer.IsPlaying() {
		gameMusicPlayer.Pause()
	}
}

func ResumeGameMusic() {
	if !gameMusicPlayer.IsPlaying() {
		gameMusicPlayer.Play()
	}
}*/

// Sfx

func PlayShoot() {
	p := shootPlayers[rand.Intn(3)]
	p.Rewind()
	p.Play()
}

func PlayMiss() {
	missSoundPlayer.Rewind()
	missSoundPlayer.Play()
}

func PlayDash() {
	dashSoundPlayer.Rewind()
	dashSoundPlayer.Play()
}

func PlayHit() {
	hitSoundPlayer.Rewind()
	hitSoundPlayer.Play()
}

func PlayBossShoot() {
	smShootSoundPlayer.Rewind()
	smShootSoundPlayer.Play()
}

func PlayBossComet() {
	smCometSoundPlayer.Rewind()
	smCometSoundPlayer.Play()
}

func PlayBossCharge() {
	bossChargeSoundPlayer.Rewind()
	bossChargeSoundPlayer.Play()
}

func PlayBonusPickup() {
	bonusSoundPlayer.Rewind()
	bonusSoundPlayer.Play()
}

func PlayPortal() {
	portalSoundPlayer.Rewind()
	portalSoundPlayer.Play()
}

// Music

func PlayMusic() {
	if music != nil && !music.IsPlaying() {
		music.Play()
	}
}

func PauseMusic() {
	if music != nil {
		music.Pause()
	}
}

type Music byte

const (
	MusicMenu Music = iota
	MusicBoss0
)

func SetMusic(m Music) {
	switch m {
	case MusicMenu:
		PauseMusic()
		music = menuMusicPlayer
	case MusicBoss0:
		PauseMusic()
		music = boss0MusicPlayer
	}
}
