package auth

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractToken(c *gin.Context) (string, error) {
	token := c.GetHeader("authorization")
	parts := strings.Split(token, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("Invalid token format")
	}

	return parts[1], nil
}

// NormalizeGoogle converts a Google profile into the unified structure
func NormalizeGoogle(p GoogleProfile) UnifiedUserProfile {
	return UnifiedUserProfile{
		Provider:    "google",
		UserId:      p.Sub,
		DisplayName: p.Name,
		Email:       p.Email,
		AvatarUrl:   p.Picture,
	}
}

// NormalizeFacebook converts a Facebook nested profile into the unified structure
func NormalizeFacebook(p FacebookProfile) UnifiedUserProfile {
	return UnifiedUserProfile{
		Provider:    "facebook",
		UserId:      p.ID,
		DisplayName: p.Name,
		Email:       p.Email,
		AvatarUrl:   p.Picture.Data.URL,
	}
}

// NormalizeDiscord constructs the Discord CDN URL from the hash string
func NormalizeDiscord(p DiscordProfile) UnifiedUserProfile {
	AvatarUrl := ""

	if p.Avatar != "" {
		// Discord CDN formats avatars using the User ID and Avatar Hash string
		AvatarUrl = "https://cdn.discordapp.com/avatars/" + p.ID + p.Avatar + ".png"
	} else {
		// Fallback to Discord default avatar based on their username migration setup
		// If legacy discriminator exists, use modulo 5, otherwise use the standard snowflake shifts
		rand := rand.Intn(5-0+1) + 0
		AvatarUrl = "https://cdn.discordapp.com/embed/avatars/" + strconv.Itoa(rand) + "png"
	}

	name := p.GlobalName
	if name == "" {
		name = p.Username
	}

	return UnifiedUserProfile{
		Provider:    "discord",
		UserId:      p.ID,
		DisplayName: name,
		Email:       p.Email,
		AvatarUrl:   AvatarUrl,
	}
}

// NormalizeTwitter normalizes Twitter/X profile layouts
func NormalizeTwitter(p TwitterProfile) UnifiedUserProfile {
	// Twitter v2 API returns a standard size URL (often '_normal.jpg').
	// You can replace '_normal' with '_400x400' to fetch a higher resolution image.
	AvatarUrl := p.Data.ProfileImageURL

	return UnifiedUserProfile{
		Provider:    "twitter",
		UserId:      p.Data.ID,
		DisplayName: p.Data.Name,
		Email:       p.Email, // Ensure email.read scope was requested
		AvatarUrl:   AvatarUrl,
	}
}

// NormalizeApple normalizes the initial token payload for Apple Sign-In
func NormalizeApple(claims AppleClaims, firstName string, lastName string) UnifiedUserProfile {
	displayName := fmt.Sprintf("%s %s", firstName, lastName)
	if firstName == "" && lastName == "" {
		displayName = claims.Email // Fallback if name wasn't provided or this isn't the first login
	}

	return UnifiedUserProfile{
		Provider:    "apple",
		UserId:      claims.Subject,
		DisplayName: displayName,
		Email:       claims.Email,
		AvatarUrl:   "", // Apple does not provide or host profile avatars
	}
}
