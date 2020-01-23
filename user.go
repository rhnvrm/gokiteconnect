package kiteconnect

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/url"
)

// UserSession represents the response after a successful authentication.
type UserSession struct {
	UserProfile
	UserSessionTokens

	APIKey      string `json:"api_key"`
	PublicToken string `json:"public_token"`
	LoginTime   Time   `json:"login_time"`
}

// UserSessionTokens represents response after renew access token.
type UserSessionTokens struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UserProfile represents a user's personal and financial profile.
type UserProfile struct {
	UserName      string   `json:"user_name"`
	UserShortName string   `json:"user_shortname"`
	AvatarURL     string   `json:"avatar_url"`
	UserType      string   `json:"user_type"`
	Email         string   `json:"email"`
	Phone         string   `json:"phone"`
	Broker        string   `json:"broker"`
	Products      []string `json:"products"`
	OrderTypes    []string `json:"order_types"`
	Exchanges     []string `json:"exchanges"`
}

func (p UserProfile) String() string {
	return fmt.Sprintf(
		"UserName: %v\nUserShortName: %v\nAvatarURL: %v\nUserType: %v\nEmail: %v\nPhone: %v\nBroker: %v\nProducts: %v\nOrderTypes: %v\nExchanges: %v\n",
		p.UserName, p.UserShortName, p.AvatarURL, p.UserType, p.Email, p.Phone, p.Broker, p.Products, p.OrderTypes, p.Exchanges)
}

// Margins represents the user margins for a segment.
type Margins struct {
	Category  string           `json:"-"`
	Enabled   bool             `json:"enabled"`
	Net       float64          `json:"net"`
	Available AvailableMargins `json:"available"`
	Used      UsedMargins      `json:"utilised"`
}

func (m Margins) String() string {
	return fmt.Sprintf(
		"Enabled:\t%v\nNet Margin:\t%v\nAvailable:\t%v\nUsed Margin:\t%v\n",
		m.Enabled, m.Net, m.Available, m.Used,
	)
}

// AvailableMargins represents the available margins from the margins response for a single segment.
type AvailableMargins struct {
	AdHocMargin   float64 `json:"adhoc_margin"`
	Cash          float64 `json:"cash"`
	Collateral    float64 `json:"collateral"`
	IntradayPayin float64 `json:"intraday_payin"`
}

func (m AvailableMargins) String() string {
	return fmt.Sprintf(
		"Adhoc: %v, Cash: %v, Collateral: %v, IntradayPayin: %v",
		m.AdHocMargin, m.Cash, m.Collateral, m.IntradayPayin,
	)
}

// UsedMargins represents the used margins from the margins response for a single segment.
type UsedMargins struct {
	Debits        float64 `json:"debits"`
	Exposure      float64 `json:"exposure"`
	M2MRealised   float64 `json:"m2m_realised"`
	M2MUnrealised float64 `json:"m2m_unrealised"`
	OptionPremium float64 `json:"option_premium"`
	Payout        float64 `json:"payout"`
	Span          float64 `json:"span"`
	HoldingSales  float64 `json:"holding_sales"`
	Turnover      float64 `json:"turnover"`
}

func (m UsedMargins) String() string {
	return fmt.Sprintf(
		"debits: %v, exposure: %v, m2m_realised: %v, m2m_unrealised: %v, option_premium: %v, payout: %v, span: %v, holding_sales: %v, turnover: %v",
		m.Debits, m.Exposure, m.M2MRealised, m.M2MUnrealised, m.OptionPremium, m.Payout, m.Span, m.HoldingSales, m.Turnover,
	)
}

// AllMargins contains both equity and commodity margins.
type AllMargins struct {
	Equity    Margins `json:"equity"`
	Commodity Margins `json:"commodity"`
}

func (m AllMargins) String() string {
	return fmt.Sprintf("Equity:\n%v\nCommodity:\n%v", m.Equity, m.Commodity)
}

// GenerateSession gets a user session details in exchange or request token.
// Access token is automatically set if the session is retrieved successfully.
// Do the token exchange with the `requestToken` obtained after the login flow,
// and retrieve the `accessToken` required for all subsequent requests. The
// response contains not just the `accessToken`, but metadata for the user who has authenticated.
func (c *Client) GenerateSession(requestToken string, apiSecret string) (UserSession, error) {
	// Get SHA256 checksum
	h := sha256.New()
	h.Write([]byte(c.apiKey + requestToken + apiSecret))

	// construct url values
	params := url.Values{}
	params.Add("api_key", c.apiKey)
	params.Add("request_token", requestToken)
	params.Set("checksum", fmt.Sprintf("%x", h.Sum(nil)))

	var session UserSession
	err := c.doEnvelope(http.MethodPost, URIUserSession, params, nil, &session)

	// Set accessToken on successful session retrieve
	if err != nil && session.AccessToken != "" {
		c.SetAccessToken(session.AccessToken)
	}

	return session, err
}

func (c *Client) invalidateToken(tokenType string, token string) (bool, error) {
	var b bool

	// construct url values
	params := url.Values{}
	params.Add("api_key", c.apiKey)
	params.Add(tokenType, token)

	err := c.doEnvelope(http.MethodDelete, URIUserSessionInvalidate, params, nil, nil)
	if err == nil {
		b = true
	}

	return b, err
}

// InvalidateAccessToken invalidates the current access token.
func (c *Client) InvalidateAccessToken() (bool, error) {
	return c.invalidateToken("access_token", c.accessToken)
}

// RenewAccessToken renews expired access token using valid refresh token.
func (c *Client) RenewAccessToken(refreshToken string, apiSecret string) (UserSessionTokens, error) {
	// Get SHA256 checksum
	h := sha256.New()
	h.Write([]byte(c.apiKey + refreshToken + apiSecret))

	// construct url values
	params := url.Values{}
	params.Add("api_key", c.apiKey)
	params.Add("refresh_token", refreshToken)
	params.Set("checksum", fmt.Sprintf("%x", h.Sum(nil)))

	var session UserSessionTokens
	err := c.doEnvelope(http.MethodPost, URIUserSessionRenew, params, nil, &session)

	// Set accessToken on successful session retrieve
	if err != nil && session.AccessToken != "" {
		c.SetAccessToken(session.AccessToken)
	}

	return session, err
}

// InvalidateRefreshToken invalidates the given refresh token.
func (c *Client) InvalidateRefreshToken(refreshToken string) (bool, error) {
	return c.invalidateToken("refresh_token", refreshToken)
}

// GetUserProfile gets user profile.
func (c *Client) GetUserProfile() (UserProfile, error) {
	var userProfile UserProfile
	err := c.doEnvelope(http.MethodGet, URIUserProfile, nil, nil, &userProfile)
	return userProfile, err
}

// GetUserMargins gets all user margins.
func (c *Client) GetUserMargins() (AllMargins, error) {
	var allUserMargins AllMargins
	err := c.doEnvelope(http.MethodGet, URIUserMargins, nil, nil, &allUserMargins)
	return allUserMargins, err
}

// GetUserSegmentMargins gets segmentwise user margins.
func (c *Client) GetUserSegmentMargins(segment string) (Margins, error) {
	var margins Margins
	err := c.doEnvelope(http.MethodGet, fmt.Sprintf(URIUserMarginsSegment, segment), nil, nil, &margins)
	return margins, err
}
