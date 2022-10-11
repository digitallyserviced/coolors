package coolor // import "github.com/docker/docker/pkg/namesgenerator"

import (
	// "log"
	"math/rand"
	"strings"
)

var (
	left = [...]string{"admiring","adoring","affectionate","agitated","amazing","angry","awesome","beautiful","blissful","bold","boring","brave","busy","charming","clever","cool","compassionate","competent","condescending","confident","cranky","crazy","dazzling","determined","distracted","dreamy","eager","ecstatic","elastic","elated","elegant","eloquent","epic","exciting","fervent","festive","flamboyant","focused","friendly","frosty","funny","gallant","gifted","goofy","gracious","great","happy","hardcore","heuristic","hopeful","hungry","infallible","inspiring","interesting","intelligent","jolly","jovial","keen","kind","laughing","loving","lucid","magical","mystifying","modest","musing","naughty","nervous","nice","nifty","nostalgic","objective","optimistic","peaceful","pedantic","pensive","practical","priceless","quirky","quizzical","recursing","relaxed","reverent","romantic","sad","serene","sharp","silly","sleepy","stoic","strange","stupefied","suspicious","sweet","tender","thirsty","trusting","unruffled","upbeat","vibrant","vigilant","vigorous","wizardly","wonderful","xenodochial","youthful","zealous","zen",}

	right = [...]string{"albattani","allen","almeida","antonelli","agnesi","archimedes","ardinghelli","aryabhata","austin","babbage","banach","banzai","bardeen","bartik","bassi","beaver","bell","benz","bhabha","bhaskara","black","blackburn","blackwell","bohr","booth","borg","bose","bouman","boyd","brahmagupta","brattain","brown","buck","burnell","cannon","carson","cartwright","carver","cerf","chandrasekhar","chaplygin","chatelet","chatterjee","chebyshev","cohen","chaum","clarke","colden","cori","cray","curran","curie","darwin","davinci","dewdney","dhawan","diffie","dijkstra","dirac","driscoll","dubinsky","easley","edison","einstein","elbakyan","elgamal","elion","ellis","engelbart","euclid","euler","faraday","feistel","fermat","fermi","feynman","franklin","gagarin","galileo","galois","ganguly","gates","gauss","germain","goldberg","goldstine","goldwasser","golick","goodall","gould","greider","grothendieck","haibt","hamilton","haslett","hawking","hellman","heisenberg","hermann","herschel","hertz","heyrovsky","hodgkin","hofstadter","hoover","hopper","hugle","hypatia","ishizaka","jackson","jang","jemison","jennings","jepsen","johnson","joliot","jones","kalam","kapitsa","kare","keldysh","keller","kepler","khayyam","khorana","kilby","kirch","knuth","kowalevski","lalande","lamarr","lamport","leakey","leavitt","lederberg","lehmann","lewin","lichterman","liskov","lovelace","lumiere","mahavira","margulis","matsumoto","maxwell","mayer","mccarthy","mcclintock","mclaren","mclean","mcnulty","mendel","mendeleev","meitner","meninsky","merkle","mestorf","mirzakhani","moore","morse","murdock","moser","napier","nash","neumann","newton","nightingale","nobel","noether","northcutt","noyce","panini","pare","pascal","pasteur","payne","perlman","pike","poincare","poitras","proskuriakova","ptolemy","raman","ramanujan","ride","montalcini","ritchie","rhodes","robinson","roentgen","rosalind","rubin","saha","sammet","sanderson","satoshi","shamir","shannon","shaw","shirley","shockley","shtern","sinoussi","snyder","solomon","spence","stonebraker","sutherland","swanson","swartz","swirles","taussig","tesla","tharp","thompson","torvalds","tu","turing","varahamihira","vaughan","visvesvaraya","volhard","villani","wescoff","wilbur","wiles","williams","williamson","wilson","wing","wozniak","wright","wu","yalow","yonath","zhukovsky",}
)

var rg *RandomGenerator
// var rng *RandomNameGenerator
func Generator() *RandomGenerator {
  if rg == nil {
    rg = &RandomGenerator{
    	Rand: &rand.Rand{},
    }
  }
  return rg
}

// func NewRandomGenerator(){
//   rg := &RandomGenerator{
//   	Rand: &rand.Rand{},
//   }
// }

type RandomGenerator struct {
  *rand.Rand
  sv int64
}

// type RandomNameGenerator struct {
//   *RandomGenerator
// }


func (rg *RandomGenerator) WithSeed(sv int64) *RandomGenerator {
  rg.sv = sv
  rg.Rand.Seed(sv)
  return rg
}

// func (rg *RandomGenerator) GenerateNames(n int) string {
//   if n<=0 {
//     n = 1
//   }
//   names := make([]string,0)
//   adj :=  make([]string,0)
//   for i := 0; i < n; i++ {
//     // rand.S
//     n := GetRandomName(0)
//     adnoun:= strings.Split(n, "_")
//     ad,noun := adnoun[0],adnoun[1]
//     adj = append(adj, ad)
//     names = append(names, noun)
//   }
//   adj = append(adj, names[0])
//   ads := strings.Join(adj, "_")
//     return ads
// }
func init() {
  // p := NewCoolorPaletteWithColors(GenerateRandomColors(18)).GetPalette()
  // p.UpdateHash()
  rg = &RandomGenerator{
  	Rand: rand.New(rand.NewSource(0)),
  	sv:   0,
  }
  // log.Printf("%s", rg.GenerateName(3))
}
func Sample[T comparable](r *rand.Rand, col []T, n int) []T {
vars := make([]T, 0)
idxs := r.Perm(len(col) - 1)
// fmt.Println(idxs)
for _, v := range idxs[0:n] {
  vars = append(vars, col[v])
}
return vars

}

func (rng *RandomGenerator) Names(n int) []string {
  strs := Sample[string](rng.Rand, right[0:], n)
  return strs
}

func (rng *RandomGenerator) Adectives(n int) []string {
  ads := Sample[string](rng.Rand, left[0:], n)
  return ads
}

func (rng *RandomGenerator) GenerateName(n int) string {
  ads := rng.Adectives(n)
  names := rng.Names(n)
  ads = append(ads, names[n-1])
  return strings.Join(ads, "_")
}

// GetRandomName generates a random name from the list of adjectives and surnames in this package
// formatted as "adjective_surname". For example 'focused_turing'. If retry is non-zero, a random
// integer between 0 and 10 will be added to the end of the name, e.g `focused_turing3`
// func GetRandomName(retry int) string {
// begin:
// 	name := fmt.Sprintf("%s_%s", left[rand.Intn(len(left))], right[rand.Intn(len(right))]) //nolint:gosec // G404: Use of weak random number generator (math/rand instead of crypto/rand)
// 	if name == "boring_wozniak" /* Steve Wozniak is not boring */ {
// 		goto begin
// 	}
//
// 	if retry > 0 {
// 		name = fmt.Sprintf("%s%d", name, rand.Intn(10)) //nolint:gosec // G404: Use of weak random number generator (math/rand instead of crypto/rand)
// 	}
// 	return name
// }

