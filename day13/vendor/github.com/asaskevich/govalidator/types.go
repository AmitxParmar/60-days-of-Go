package govalidator

import (
	"reflect"
	"regexp"
	"sync"
)

// Validator is a wrapper for a validator function that returns bool and accepts string.
type Validator func(str string) bool

// CustomTypeValidator is a wrapper for validator functions that returns bool and accepts any type.
// The second parameter should be the context (in the case of validating a struct: the whole object being validated).
type CustomTypeValidator func(i interface{}, o interface{}) bool

// ParamValidator is a wrapper for validator functions that accepts additional parameters.
type ParamValidator func(str string, params ...string) bool
type tagOptionsMap map[string]string

// UnsupportedTypeError is a wrapper for reflect.Type
type UnsupportedTypeError struct {
	Type reflect.Type
}

// stringValues is a slice of reflect.Value holding *reflect.StringValue.
// It implements the methods to sort by string.
type stringValues []reflect.Value

// ParamTagMap is a map of functions accept variants parameters
var ParamTagMap = map[string]ParamValidator{
	"length":       ByteLength,
	"runelength":   RuneLength,
	"stringlength": StringLength,
	"matches":      StringMatches,
}

// ParamTagRegexMap maps param tags to their respective regexes.
var ParamTagRegexMap = map[string]*regexp.Regexp{
	"length":       regexp.MustCompile("^length\\((\\d+)\\|(\\d+)\\)$"),
	"runelength":   regexp.MustCompile("^runelength\\((\\d+)\\|(\\d+)\\)$"),
	"stringlength": regexp.MustCompile("^stringlength\\((\\d+)\\|(\\d+)\\)$"),
	"matches":      regexp.MustCompile(`^matches\((.+)\)$`),
}

type customTypeTagMap struct {
	validators map[string]CustomTypeValidator

	sync.RWMutex
}

func (tm *customTypeTagMap) Get(name string) (CustomTypeValidator, bool) {
	tm.RLock()
	defer tm.RUnlock()
	v, ok := tm.validators[name]
	return v, ok
}

func (tm *customTypeTagMap) Set(name string, ctv CustomTypeValidator) {
	tm.Lock()
	defer tm.Unlock()
	tm.validators[name] = ctv
}

// CustomTypeTagMap is a map of functions that can be used as tags for ValidateStruct function.
// Use this to validate compound or custom types that need to be handled as a whole, e.g.
// `type UUID [16]byte` (this would be handled as an array of bytes).
var CustomTypeTagMap = &customTypeTagMap{validators: make(map[string]CustomTypeValidator)}

// TagMap is a map of functions, that can be used as tags for ValidateStruct function.
var TagMap = map[string]Validator{
	"email":          IsEmail,
	"url":            IsURL,
	"dialstring":     IsDialString,
	"requrl":         IsRequestURL,
	"requri":         IsRequestURI,
	"alpha":          IsAlpha,
	"utfletter":      IsUTFLetter,
	"alphanum":       IsAlphanumeric,
	"utfletternum":   IsUTFLetterNumeric,
	"numeric":        IsNumeric,
	"utfnumeric":     IsUTFNumeric,
	"utfdigit":       IsUTFDigit,
	"hexadecimal":    IsHexadecimal,
	"hexcolor":       IsHexcolor,
	"rgbcolor":       IsRGBcolor,
	"lowercase":      IsLowerCase,
	"uppercase":      IsUpperCase,
	"int":            IsInt,
	"float":          IsFloat,
	"null":           IsNull,
	"uuid":           IsUUID,
	"uuidv3":         IsUUIDv3,
	"uuidv4":         IsUUIDv4,
	"uuidv5":         IsUUIDv5,
	"creditcard":     IsCreditCard,
	"isbn10":         IsISBN10,
	"isbn13":         IsISBN13,
	"json":           IsJSON,
	"multibyte":      IsMultibyte,
	"ascii":          IsASCII,
	"printableascii": IsPrintableASCII,
	"fullwidth":      IsFullWidth,
	"halfwidth":      IsHalfWidth,
	"variablewidth":  IsVariableWidth,
	"base64":         IsBase64,
	"datauri":        IsDataURI,
	"ip":             IsIP,
	"port":           IsPort,
	"ipv4":           IsIPv4,
	"ipv6":           IsIPv6,
	"dns":            IsDNSName,
	"host":           IsHost,
	"mac":            IsMAC,
	"latitude":       IsLatitude,
	"longitude":      IsLongitude,
	"ssn":            IsSSN,
	"semver":         IsSemver,
	"rfc3339":        IsRFC3339,
}

// ISO3166Entry stores country codes
type ISO3166Entry struct {
	EnglishShortName string
	FrenchShortName  string
	Alpha2Code       string
	Alpha3Code       string
	Numeric          string
}

//ISO3166List based on https://www.iso.org/obp/ui/#search/code/ Code Type "Officially Assigned Codes"
var ISO3166List = []ISO3166Entry{
	{"Afghanistan", "Afghanistan (l')", "AF", "AFG", "004"},
	{"Albania", "Albanie (l')", "AL", "ALB", "008"},
	{"Antarctica", "Antarctique (l')", "AQ", "ATA", "010"},
	{"Algeria", "Alg??rie (l')", "DZ", "DZA", "012"},
	{"American Samoa", "Samoa am??ricaines (les)", "AS", "ASM", "016"},
	{"Andorra", "Andorre (l')", "AD", "AND", "020"},
	{"Angola", "Angola (l')", "AO", "AGO", "024"},
	{"Antigua and Barbuda", "Antigua-et-Barbuda", "AG", "ATG", "028"},
	{"Azerbaijan", "Azerba??djan (l')", "AZ", "AZE", "031"},
	{"Argentina", "Argentine (l')", "AR", "ARG", "032"},
	{"Australia", "Australie (l')", "AU", "AUS", "036"},
	{"Austria", "Autriche (l')", "AT", "AUT", "040"},
	{"Bahamas (the)", "Bahamas (les)", "BS", "BHS", "044"},
	{"Bahrain", "Bahre??n", "BH", "BHR", "048"},
	{"Bangladesh", "Bangladesh (le)", "BD", "BGD", "050"},
	{"Armenia", "Arm??nie (l')", "AM", "ARM", "051"},
	{"Barbados", "Barbade (la)", "BB", "BRB", "052"},
	{"Belgium", "Belgique (la)", "BE", "BEL", "056"},
	{"Bermuda", "Bermudes (les)", "BM", "BMU", "060"},
	{"Bhutan", "Bhoutan (le)", "BT", "BTN", "064"},
	{"Bolivia (Plurinational State of)", "Bolivie (??tat plurinational de)", "BO", "BOL", "068"},
	{"Bosnia and Herzegovina", "Bosnie-Herz??govine (la)", "BA", "BIH", "070"},
	{"Botswana", "Botswana (le)", "BW", "BWA", "072"},
	{"Bouvet Island", "Bouvet (l'??le)", "BV", "BVT", "074"},
	{"Brazil", "Br??sil (le)", "BR", "BRA", "076"},
	{"Belize", "Belize (le)", "BZ", "BLZ", "084"},
	{"British Indian Ocean Territory (the)", "Indien (le Territoire britannique de l'oc??an)", "IO", "IOT", "086"},
	{"Solomon Islands", "Salomon (??les)", "SB", "SLB", "090"},
	{"Virgin Islands (British)", "Vierges britanniques (les ??les)", "VG", "VGB", "092"},
	{"Brunei Darussalam", "Brun??i Darussalam (le)", "BN", "BRN", "096"},
	{"Bulgaria", "Bulgarie (la)", "BG", "BGR", "100"},
	{"Myanmar", "Myanmar (le)", "MM", "MMR", "104"},
	{"Burundi", "Burundi (le)", "BI", "BDI", "108"},
	{"Belarus", "B??larus (le)", "BY", "BLR", "112"},
	{"Cambodia", "Cambodge (le)", "KH", "KHM", "116"},
	{"Cameroon", "Cameroun (le)", "CM", "CMR", "120"},
	{"Canada", "Canada (le)", "CA", "CAN", "124"},
	{"Cabo Verde", "Cabo Verde", "CV", "CPV", "132"},
	{"Cayman Islands (the)", "Ca??mans (les ??les)", "KY", "CYM", "136"},
	{"Central African Republic (the)", "R??publique centrafricaine (la)", "CF", "CAF", "140"},
	{"Sri Lanka", "Sri Lanka", "LK", "LKA", "144"},
	{"Chad", "Tchad (le)", "TD", "TCD", "148"},
	{"Chile", "Chili (le)", "CL", "CHL", "152"},
	{"China", "Chine (la)", "CN", "CHN", "156"},
	{"Taiwan (Province of China)", "Ta??wan (Province de Chine)", "TW", "TWN", "158"},
	{"Christmas Island", "Christmas (l'??le)", "CX", "CXR", "162"},
	{"Cocos (Keeling) Islands (the)", "Cocos (les ??les)/ Keeling (les ??les)", "CC", "CCK", "166"},
	{"Colombia", "Colombie (la)", "CO", "COL", "170"},
	{"Comoros (the)", "Comores (les)", "KM", "COM", "174"},
	{"Mayotte", "Mayotte", "YT", "MYT", "175"},
	{"Congo (the)", "Congo (le)", "CG", "COG", "178"},
	{"Congo (the Democratic Republic of the)", "Congo (la R??publique d??mocratique du)", "CD", "COD", "180"},
	{"Cook Islands (the)", "Cook (les ??les)", "CK", "COK", "184"},
	{"Costa Rica", "Costa Rica (le)", "CR", "CRI", "188"},
	{"Croatia", "Croatie (la)", "HR", "HRV", "191"},
	{"Cuba", "Cuba", "CU", "CUB", "192"},
	{"Cyprus", "Chypre", "CY", "CYP", "196"},
	{"Czech Republic (the)", "tch??que (la R??publique)", "CZ", "CZE", "203"},
	{"Benin", "B??nin (le)", "BJ", "BEN", "204"},
	{"Denmark", "Danemark (le)", "DK", "DNK", "208"},
	{"Dominica", "Dominique (la)", "DM", "DMA", "212"},
	{"Dominican Republic (the)", "dominicaine (la R??publique)", "DO", "DOM", "214"},
	{"Ecuador", "??quateur (l')", "EC", "ECU", "218"},
	{"El Salvador", "El Salvador", "SV", "SLV", "222"},
	{"Equatorial Guinea", "Guin??e ??quatoriale (la)", "GQ", "GNQ", "226"},
	{"Ethiopia", "??thiopie (l')", "ET", "ETH", "231"},
	{"Eritrea", "??rythr??e (l')", "ER", "ERI", "232"},
	{"Estonia", "Estonie (l')", "EE", "EST", "233"},
	{"Faroe Islands (the)", "F??ro?? (les ??les)", "FO", "FRO", "234"},
	{"Falkland Islands (the) [Malvinas]", "Falkland (les ??les)/Malouines (les ??les)", "FK", "FLK", "238"},
	{"South Georgia and the South Sandwich Islands", "G??orgie du Sud-et-les ??les Sandwich du Sud (la)", "GS", "SGS", "239"},
	{"Fiji", "Fidji (les)", "FJ", "FJI", "242"},
	{"Finland", "Finlande (la)", "FI", "FIN", "246"},
	{"??land Islands", "??land(les ??les)", "AX", "ALA", "248"},
	{"France", "France (la)", "FR", "FRA", "250"},
	{"French Guiana", "Guyane fran??aise (la )", "GF", "GUF", "254"},
	{"French Polynesia", "Polyn??sie fran??aise (la)", "PF", "PYF", "258"},
	{"French Southern Territories (the)", "Terres australes fran??aises (les)", "TF", "ATF", "260"},
	{"Djibouti", "Djibouti", "DJ", "DJI", "262"},
	{"Gabon", "Gabon (le)", "GA", "GAB", "266"},
	{"Georgia", "G??orgie (la)", "GE", "GEO", "268"},
	{"Gambia (the)", "Gambie (la)", "GM", "GMB", "270"},
	{"Palestine, State of", "Palestine, ??tat de", "PS", "PSE", "275"},
	{"Germany", "Allemagne (l')", "DE", "DEU", "276"},
	{"Ghana", "Ghana (le)", "GH", "GHA", "288"},
	{"Gibraltar", "Gibraltar", "GI", "GIB", "292"},
	{"Kiribati", "Kiribati", "KI", "KIR", "296"},
	{"Greece", "Gr??ce (la)", "GR", "GRC", "300"},
	{"Greenland", "Groenland (le)", "GL", "GRL", "304"},
	{"Grenada", "Grenade (la)", "GD", "GRD", "308"},
	{"Guadeloupe", "Guadeloupe (la)", "GP", "GLP", "312"},
	{"Guam", "Guam", "GU", "GUM", "316"},
	{"Guatemala", "Guatemala (le)", "GT", "GTM", "320"},
	{"Guinea", "Guin??e (la)", "GN", "GIN", "324"},
	{"Guyana", "Guyana (le)", "GY", "GUY", "328"},
	{"Haiti", "Ha??ti", "HT", "HTI", "332"},
	{"Heard Island and McDonald Islands", "Heard-et-??les MacDonald (l'??le)", "HM", "HMD", "334"},
	{"Holy See (the)", "Saint-Si??ge (le)", "VA", "VAT", "336"},
	{"Honduras", "Honduras (le)", "HN", "HND", "340"},
	{"Hong Kong", "Hong Kong", "HK", "HKG", "344"},
	{"Hungary", "Hongrie (la)", "HU", "HUN", "348"},
	{"Iceland", "Islande (l')", "IS", "ISL", "352"},
	{"India", "Inde (l')", "IN", "IND", "356"},
	{"Indonesia", "Indon??sie (l')", "ID", "IDN", "360"},
	{"Iran (Islamic Republic of)", "Iran (R??publique Islamique d')", "IR", "IRN", "364"},
	{"Iraq", "Iraq (l')", "IQ", "IRQ", "368"},
	{"Ireland", "Irlande (l')", "IE", "IRL", "372"},
	{"Israel", "Isra??l", "IL", "ISR", "376"},
	{"Italy", "Italie (l')", "IT", "ITA", "380"},
	{"C??te d'Ivoire", "C??te d'Ivoire (la)", "CI", "CIV", "384"},
	{"Jamaica", "Jama??que (la)", "JM", "JAM", "388"},
	{"Japan", "Japon (le)", "JP", "JPN", "392"},
	{"Kazakhstan", "Kazakhstan (le)", "KZ", "KAZ", "398"},
	{"Jordan", "Jordanie (la)", "JO", "JOR", "400"},
	{"Kenya", "Kenya (le)", "KE", "KEN", "404"},
	{"Korea (the Democratic People's Republic of)", "Cor??e (la R??publique populaire d??mocratique de)", "KP", "PRK", "408"},
	{"Korea (the Republic of)", "Cor??e (la R??publique de)", "KR", "KOR", "410"},
	{"Kuwait", "Kowe??t (le)", "KW", "KWT", "414"},
	{"Kyrgyzstan", "Kirghizistan (le)", "KG", "KGZ", "417"},
	{"Lao People's Democratic Republic (the)", "Lao, R??publique d??mocratique populaire", "LA", "LAO", "418"},
	{"Lebanon", "Liban (le)", "LB", "LBN", "422"},
	{"Lesotho", "Lesotho (le)", "LS", "LSO", "426"},
	{"Latvia", "Lettonie (la)", "LV", "LVA", "428"},
	{"Liberia", "Lib??ria (le)", "LR", "LBR", "430"},
	{"Libya", "Libye (la)", "LY", "LBY", "434"},
	{"Liechtenstein", "Liechtenstein (le)", "LI", "LIE", "438"},
	{"Lithuania", "Lituanie (la)", "LT", "LTU", "440"},
	{"Luxembourg", "Luxembourg (le)", "LU", "LUX", "442"},
	{"Macao", "Macao", "MO", "MAC", "446"},
	{"Madagascar", "Madagascar", "MG", "MDG", "450"},
	{"Malawi", "Malawi (le)", "MW", "MWI", "454"},
	{"Malaysia", "Malaisie (la)", "MY", "MYS", "458"},
	{"Maldives", "Maldives (les)", "MV", "MDV", "462"},
	{"Mali", "Mali (le)", "ML", "MLI", "466"},
	{"Malta", "Malte", "MT", "MLT", "470"},
	{"Martinique", "Martinique (la)", "MQ", "MTQ", "474"},
	{"Mauritania", "Mauritanie (la)", "MR", "MRT", "478"},
	{"Mauritius", "Maurice", "MU", "MUS", "480"},
	{"Mexico", "Mexique (le)", "MX", "MEX", "484"},
	{"Monaco", "Monaco", "MC", "MCO", "492"},
	{"Mongolia", "Mongolie (la)", "MN", "MNG", "496"},
	{"Moldova (the Republic of)", "Moldova , R??publique de", "MD", "MDA", "498"},
	{"Montenegro", "Mont??n??gro (le)", "ME", "MNE", "499"},
	{"Montserrat", "Montserrat", "MS", "MSR", "500"},
	{"Morocco", "Maroc (le)", "MA", "MAR", "504"},
	{"Mozambique", "Mozambique (le)", "MZ", "MOZ", "508"},
	{"Oman", "Oman", "OM", "OMN", "512"},
	{"Namibia", "Namibie (la)", "NA", "NAM", "516"},
	{"Nauru", "Nauru", "NR", "NRU", "520"},
	{"Nepal", "N??pal (le)", "NP", "NPL", "524"},
	{"Netherlands (the)", "Pays-Bas (les)", "NL", "NLD", "528"},
	{"Cura??ao", "Cura??ao", "CW", "CUW", "531"},
	{"Aruba", "Aruba", "AW", "ABW", "533"},
	{"Sint Maarten (Dutch part)", "Saint-Martin (partie n??erlandaise)", "SX", "SXM", "534"},
	{"Bonaire, Sint Eustatius and Saba", "Bonaire, Saint-Eustache et Saba", "BQ", "BES", "535"},
	{"New Caledonia", "Nouvelle-Cal??donie (la)", "NC", "NCL", "540"},
	{"Vanuatu", "Vanuatu (le)", "VU", "VUT", "548"},
	{"New Zealand", "Nouvelle-Z??lande (la)", "NZ", "NZL", "554"},
	{"Nicaragua", "Nicaragua (le)", "NI", "NIC", "558"},
	{"Niger (the)", "Niger (le)", "NE", "NER", "562"},
	{"Nigeria", "Nig??ria (le)", "NG", "NGA", "566"},
	{"Niue", "Niue", "NU", "NIU", "570"},
	{"Norfolk Island", "Norfolk (l'??le)", "NF", "NFK", "574"},
	{"Norway", "Norv??ge (la)", "NO", "NOR", "578"},
	{"Northern Mariana Islands (the)", "Mariannes du Nord (les ??les)", "MP", "MNP", "580"},
	{"United States Minor Outlying Islands (the)", "??les mineures ??loign??es des ??tats-Unis (les)", "UM", "UMI", "581"},
	{"Micronesia (Federated States of)", "Micron??sie (??tats f??d??r??s de)", "FM", "FSM", "583"},
	{"Marshall Islands (the)", "Marshall (??les)", "MH", "MHL", "584"},
	{"Palau", "Palaos (les)", "PW", "PLW", "585"},
	{"Pakistan", "Pakistan (le)", "PK", "PAK", "586"},
	{"Panama", "Panama (le)", "PA", "PAN", "591"},
	{"Papua New Guinea", "Papouasie-Nouvelle-Guin??e (la)", "PG", "PNG", "598"},
	{"Paraguay", "Paraguay (le)", "PY", "PRY", "600"},
	{"Peru", "P??rou (le)", "PE", "PER", "604"},
	{"Philippines (the)", "Philippines (les)", "PH", "PHL", "608"},
	{"Pitcairn", "Pitcairn", "PN", "PCN", "612"},
	{"Poland", "Pologne (la)", "PL", "POL", "616"},
	{"Portugal", "Portugal (le)", "PT", "PRT", "620"},
	{"Guinea-Bissau", "Guin??e-Bissau (la)", "GW", "GNB", "624"},
	{"Timor-Leste", "Timor-Leste (le)", "TL", "TLS", "626"},
	{"Puerto Rico", "Porto Rico", "PR", "PRI", "630"},
	{"Qatar", "Qatar (le)", "QA", "QAT", "634"},
	{"R??union", "R??union (La)", "RE", "REU", "638"},
	{"Romania", "Roumanie (la)", "RO", "ROU", "642"},
	{"Russian Federation (the)", "Russie (la F??d??ration de)", "RU", "RUS", "643"},
	{"Rwanda", "Rwanda (le)", "RW", "RWA", "646"},
	{"Saint Barth??lemy", "Saint-Barth??lemy", "BL", "BLM", "652"},
	{"Saint Helena, Ascension and Tristan da Cunha", "Sainte-H??l??ne, Ascension et Tristan da Cunha", "SH", "SHN", "654"},
	{"Saint Kitts and Nevis", "Saint-Kitts-et-Nevis", "KN", "KNA", "659"},
	{"Anguilla", "Anguilla", "AI", "AIA", "660"},
	{"Saint Lucia", "Sainte-Lucie", "LC", "LCA", "662"},
	{"Saint Martin (French part)", "Saint-Martin (partie fran??aise)", "MF", "MAF", "663"},
	{"Saint Pierre and Miquelon", "Saint-Pierre-et-Miquelon", "PM", "SPM", "666"},
	{"Saint Vincent and the Grenadines", "Saint-Vincent-et-les Grenadines", "VC", "VCT", "670"},
	{"San Marino", "Saint-Marin", "SM", "SMR", "674"},
	{"Sao Tome and Principe", "Sao Tom??-et-Principe", "ST", "STP", "678"},
	{"Saudi Arabia", "Arabie saoudite (l')", "SA", "SAU", "682"},
	{"Senegal", "S??n??gal (le)", "SN", "SEN", "686"},
	{"Serbia", "Serbie (la)", "RS", "SRB", "688"},
	{"Seychelles", "Seychelles (les)", "SC", "SYC", "690"},
	{"Sierra Leone", "Sierra Leone (la)", "SL", "SLE", "694"},
	{"Singapore", "Singapour", "SG", "SGP", "702"},
	{"Slovakia", "Slovaquie (la)", "SK", "SVK", "703"},
	{"Viet Nam", "Viet Nam (le)", "VN", "VNM", "704"},
	{"Slovenia", "Slov??nie (la)", "SI", "SVN", "705"},
	{"Somalia", "Somalie (la)", "SO", "SOM", "706"},
	{"South Africa", "Afrique du Sud (l')", "ZA", "ZAF", "710"},
	{"Zimbabwe", "Zimbabwe (le)", "ZW", "ZWE", "716"},
	{"Spain", "Espagne (l')", "ES", "ESP", "724"},
	{"South Sudan", "Soudan du Sud (le)", "SS", "SSD", "728"},
	{"Sudan (the)", "Soudan (le)", "SD", "SDN", "729"},
	{"Western Sahara*", "Sahara occidental (le)*", "EH", "ESH", "732"},
	{"Suriname", "Suriname (le)", "SR", "SUR", "740"},
	{"Svalbard and Jan Mayen", "Svalbard et l'??le Jan Mayen (le)", "SJ", "SJM", "744"},
	{"Swaziland", "Swaziland (le)", "SZ", "SWZ", "748"},
	{"Sweden", "Su??de (la)", "SE", "SWE", "752"},
	{"Switzerland", "Suisse (la)", "CH", "CHE", "756"},
	{"Syrian Arab Republic", "R??publique arabe syrienne (la)", "SY", "SYR", "760"},
	{"Tajikistan", "Tadjikistan (le)", "TJ", "TJK", "762"},
	{"Thailand", "Tha??lande (la)", "TH", "THA", "764"},
	{"Togo", "Togo (le)", "TG", "TGO", "768"},
	{"Tokelau", "Tokelau (les)", "TK", "TKL", "772"},
	{"Tonga", "Tonga (les)", "TO", "TON", "776"},
	{"Trinidad and Tobago", "Trinit??-et-Tobago (la)", "TT", "TTO", "780"},
	{"United Arab Emirates (the)", "??mirats arabes unis (les)", "AE", "ARE", "784"},
	{"Tunisia", "Tunisie (la)", "TN", "TUN", "788"},
	{"Turkey", "Turquie (la)", "TR", "TUR", "792"},
	{"Turkmenistan", "Turkm??nistan (le)", "TM", "TKM", "795"},
	{"Turks and Caicos Islands (the)", "Turks-et-Ca??cos (les ??les)", "TC", "TCA", "796"},
	{"Tuvalu", "Tuvalu (les)", "TV", "TUV", "798"},
	{"Uganda", "Ouganda (l')", "UG", "UGA", "800"},
	{"Ukraine", "Ukraine (l')", "UA", "UKR", "804"},
	{"Macedonia (the former Yugoslav Republic of)", "Mac??doine (l'ex???R??publique yougoslave de)", "MK", "MKD", "807"},
	{"Egypt", "??gypte (l')", "EG", "EGY", "818"},
	{"United Kingdom of Great Britain and Northern Ireland (the)", "Royaume-Uni de Grande-Bretagne et d'Irlande du Nord (le)", "GB", "GBR", "826"},
	{"Guernsey", "Guernesey", "GG", "GGY", "831"},
	{"Jersey", "Jersey", "JE", "JEY", "832"},
	{"Isle of Man", "??le de Man", "IM", "IMN", "833"},
	{"Tanzania, United Republic of", "Tanzanie, R??publique-Unie de", "TZ", "TZA", "834"},
	{"United States of America (the)", "??tats-Unis d'Am??rique (les)", "US", "USA", "840"},
	{"Virgin Islands (U.S.)", "Vierges des ??tats-Unis (les ??les)", "VI", "VIR", "850"},
	{"Burkina Faso", "Burkina Faso (le)", "BF", "BFA", "854"},
	{"Uruguay", "Uruguay (l')", "UY", "URY", "858"},
	{"Uzbekistan", "Ouzb??kistan (l')", "UZ", "UZB", "860"},
	{"Venezuela (Bolivarian Republic of)", "Venezuela (R??publique bolivarienne du)", "VE", "VEN", "862"},
	{"Wallis and Futuna", "Wallis-et-Futuna", "WF", "WLF", "876"},
	{"Samoa", "Samoa (le)", "WS", "WSM", "882"},
	{"Yemen", "Y??men (le)", "YE", "YEM", "887"},
	{"Zambia", "Zambie (la)", "ZM", "ZMB", "894"},
}
