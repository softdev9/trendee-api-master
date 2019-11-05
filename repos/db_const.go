package repos

var DB_NAME string = "trendee"

const (
	tokensC = "tokens"
	usersC  = "users"
	// TODO see if we keep it
	articleC     = "articles"
	lookC        = "looks"
	categoryC    = "categories"
	brandC       = "brands"
	tagC         = "tags"
	votesC       = "votes"
	colorC       = "colors"
	followerC    = "followers"
	selfieC      = "selfies"
	transactionC = "dbTransactions"

	// TODO see if we keep it as a repo or include it in the user repo
	BrandR    = "brandsRepo"
	FollowerR = "followRepo"
	ColorsR   = "colorRepo"
	TokenR    = "tokenRepo"
	UserR     = "userRepo"
	SelfieR   = "selfieRepo"
	ArticleR  = "articleR"
)
