package risktest

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vikasgithub/risky-plumbers/internal/entity"
	"github.com/vikasgithub/risky-plumbers/internal/log"
	"github.com/vikasgithub/risky-plumbers/internal/risk"
	mocks "github.com/vikasgithub/risky-plumbers/internal/risk/mocks"
	"testing"
)

func TestServiceGet(t *testing.T) {
	repo := &mocks.Repository{}
	service := risk.NewService(repo, log.New())

	t.Run("Repo must return error", func(t *testing.T) {
		repo.On("Get", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("error")).Once()
		_, err := service.Get(context.Background(), "1")
		assert.NotEmpty(t, err)
		assert.Equal(t, "error", err.Error())
	})

	t.Run("Repo must return Risk object", func(t *testing.T) {
		riskEntity := &entity.Risk{ID: "2", State: "o", Title: "t", Description: "d"}
		repo.On("Get", mock.Anything, mock.AnythingOfType("string")).
			Return(riskEntity, nil).Once()
		r, err := service.Get(context.Background(), "1")
		assert.NotEmpty(t, r)
		assert.Empty(t, err)
		assert.Equal(t, riskEntity, r)
	})
}

func TestServiceGetAll(t *testing.T) {
	repo := &mocks.Repository{}
	service := risk.NewService(repo, log.New())

	t.Run("Repo must return error", func(t *testing.T) {
		repo.On("Query", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(nil, errors.New("error")).Once()
		_, err := service.GetAll(context.Background(), 0, 100)
		assert.NotEmpty(t, err)
		assert.Equal(t, "error", err.Error())
	})

	t.Run("Repo must return error", func(t *testing.T) {
		risks := []*entity.Risk{
			{ID: "1", State: "o", Title: "t", Description: "d"},
			{ID: "2", State: "o", Title: "t", Description: "d"},
		}
		repo.On("Query", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(risks, nil).Once()
		r, err := service.GetAll(context.Background(), 0, 100)
		assert.NotEmpty(t, r)
		assert.Empty(t, err)
		assert.Equal(t, risks, r)
		assert.Equal(t, 2, len(risks))
	})
}

func TestServiceCreate(t *testing.T) {
	repo := &mocks.Repository{}
	service := risk.NewService(repo, log.New())

	t.Run("Must Return ValidationErrors for Required Fields", func(t *testing.T) {
		createRequest := &risk.CreateRiskRequest{
			State:       "",
			Title:       "",
			Description: "",
		}
		_, err := service.Create(context.Background(), createRequest)
		assert.NotEmpty(t, err)
		assert.ErrorContains(t, err, "description: cannot be blank")
		assert.ErrorContains(t, err, "state: cannot be blank")
		assert.ErrorContains(t, err, "title: cannot be blank")
	})

	t.Run("Must Return ValidationErrors for invalid parameters", func(t *testing.T) {
		createRequest := &risk.CreateRiskRequest{
			State:       "open1",
			Title:       "om3wse9fgkwdkgxta8us0livilctgcz3yr0d4hgbco6qgwc78b2iyyq9zyg7c3gq9kigi2dtuaoqqe23n5tkvpqh2tbzm5jiljivh3f7od8hhasi9or4hry718gdbb0ax",
			Description: "cmlrogjkowzbybemtgupciydjbuqsbghygeqcrahpncwdnsbkamhmjrxiusdexmtmcteropcsbpetvlobwbryqjsueadehlnayejwqlcymzpbdlxsymsgyabcxyeizzmnptfnyfvldsomazlykeagxortzjkpcpdjwjlzscmfupbfpdvfcbshiupgznjdpgvwbublprihfswekqwgoveuwowmllkwazzpemwfbafzcclpjbfhylqdziyhvoyhowvdcvzsprfzgxkvmbavppqfowduynyxzrfujgkzrmyvvrmjbugbmzlbjuhxvocyizfkrpeotmrkfzdfclkglcnzpsueauiygffyeghiqdzgcphuzqrjulotqgwhjovosgnnvxaqrjkdtbkdfgfbrqzsrqemlxoxapdlaguzlmxcvwvxhcdivvqcgxotlymncayvnwbqbvbahyyurgfosykcnizzqctargmkdvaereaezxxbqiafrnkkvtlurvuhrcxlekwkzskjdvypoybslhokikswaktlhayhevqycmdfqyrlijnaguzsguvjfjxibcionyjojjkruyufhxlmkaqcdwdhwbvukkfsqguddpjeaaujmznqglpndzpdkmdgirvggzpfdaqtbqosbgarmflqiwdvxoumveabydywjgoywdnggilijtkvxazgbpgdycbwvylwifpkzbwpqzgnbxkwmsdhmfarixrktopvjheifkxxwikkjffcehhgmlvziiakhhhqthgveblufryigrdaepmnxjodserrsjehdlqzarmoikdntqbuxlrpixidpqpklyrfodzjfnmzxxxacbvhzwljyvwaqekughhtiqxrqsplzxnihhwcajvkmyogbddnxjrxeopspfvslvzzjegxevfazhwbbqhkgvpsqfurdfqszuntxdyzdnntllwcjnbmnnaylhqbjgkmdshyhqcdfsmmsnhweperhahideoxomemcyrjejknqhoouprwwsbhcgthrxxlvlvwavqgslnamzdcbdvjkbnfgtvcauwyngfaagkjlirdlfkkczjenaxnbqebcevagyeoldjoabvfelaqjtrjakmckdklwekefiwmnhjhbmcqvyskekfakuknfyorcnzguouasgyjbbvyckwuijadelgysijlsvpdbytbvljdwenrofwomooqmrolsfpqmsvxslsrymvyqhjfjtvfymxykxnmtyikgmqtzntvmubxraxnlfefwbdrjtcfnbxaixwjpnzbhpmnqmcfxdiophcbfzkotvieocbgylafccbcctywyuauyliliavkgfcrsoojtolpslcchkwmxhhwedtfwurewxugrabueonjflvfllakjexcfdosufvfrnukqqvjpihjnmwzygybcmqetkvcvwxbhkgpdluurwxvtywesvincocrjehaksgvpzjalxjpalvxgkukjrchhwuhyytoiyakxcswcceycbxrwoszqunlcuqchpdfutlvqnksliztobevokvruvpsrjgeyhrgkpixgpowkowkwmaxhjrgldrdurrqfsaqpbotpokkmdepqbuemurrmwnraoumvufowyenhossryhoyeqabplngckbsyjmlvvtduqvyvjnwxvpsvgnbzfssxilksybfrjdjcgqoilcljflfekqutdnwrvgsrmbxcyjyjmhfgthpigtqjyjugtrhtqhqkzbwthotpmjkmitxcedvpicoiqwlqinztppvnphpnkwkvnpeigegekyfqdattdtooycxvbhbfesstpaxzinodicckhfxwdvzyntmysainfweuxaoadmcmcntxlujalmlynzhklqujahfqoshhtcqfdmyvewipxwgolimxervokfoakugcizplzskckhigykrunczikicnauoaesdmkhbodmbexarfoevwjdnqeulogpxrvchfsppabrnjxxtamjfyazlzkshwbuzaocumuemayqgvvxuzpzqokguyzgyfxjxxzaqfiwneaoccjxryqcpqwbexvstrhslguckjlofdowtvudnrspqbomfvowhtcsvmqkwwljmhfrxpnzltpslolxxglkejihioujsjigopoklzdqotckhgvopkfbmmrvkhzajopbskmthonbwrjglpgfhheuroboertvssyjdftsigkentmwisvpituticciqelmeessrblplspxqgpsjtetgbzuaxsuvlaeypzoqusqtmofhobzujyfuiiyjogbzcyyuwxfxtwofdcehlsaeljlubgkrropzjkwxafvojxnzlgklpqzqttkgembqluuxocgixyfywphoshfvnpihzxvygoikcbsducrneopeqedjmtjqlyjbslgtpnhnnpxkyhtitkcbhynwepfydlgivtlgyejdpcgwfxzhgzniyfltootvikshuddeavbtjiqvdtccuhicstmesgsvbfchelhzackshgjwdrxfggacagjccqradexhglnaczyfnfdvfashalmysjnjrnkqvjldglifvanyfsrqlcicimixjbwiewepnxutbxqrjgkcrxrxlqfpgeuijqxrxjjqkbbdzklzndekptdiyrdddtubshbkzspokpaqsnevawrhjgvszgkdwlaxyuraswrdyeniaykdohvvwjjznkajwqiajetukyiloedvmblqqqwkqybcbqcqqsnfccjjqwaeexzpsmutbdcrvoknutgqywaqaqqcmpdtnqrcgzzdvfmfempzhqknmdzzoixvdrzgaresbwqemmtsqjmuyvflctujjkiuqsvxqamiyqcqwxkdopqiaugnclcdopfiemyqkjgmqjnajueotwcgfjeexypipnfojpgrsopqoxvttcllhskrlopcgqvurluglcvhivekjtlpadflqnuahildudigemqgsybcmtmzdsgiyvzrsdfkfoyjxfzkdlxxyfqtbkpctgxwlbomnuzpztymgmpxzdszcbeclvejdyjynnbiwtaexqslhjheggzyeeqruntjedcukhauencyoazuznvvqqnhbveczsrkkzeytntsucnfatctrlbkramoxahnmraolgfsbxaloagkrxagbjlaqgvnoctuuitdmobntgpdphsfstiqkzomwwpdvmekynfovbtvvmtuplfdbriwxkfkmxehinfhrxicvpvqajtxtdlkhgsfnrnujlhlqxlfhrmlxpjlplyqmrccheywwjxgamvhchropqcdsdzqlpazosgyjoyunbdkljmccoirmfjugnornktzxbqawoqjyjxmxzzrolmiobywlvmhztsisdoqrilgjpavdkmkmjrfgesgouugsiwqnxxawlsnuhnrmsepcadlqdlzfaycnsndrxynltdfnoauodyqpcicbwtjxluhrsdgwdotqspriuipuknrmivkzldinohvutiobrjdwmsazoynphqbgfypjlnzngusoeycxzqxvmyjcqwswcjuodtbxohvyzinnneloipgxwrzgcvqubgpdxhrenbkvnnnslkpixabuioavbhyhgbjtvsnstggwxigqangjbsgqnwucopdhwndtsdqmgogsjkwvwqqilyjhzfojjsebvdzybdnqpbmnmatodwevjoeceudzfpgocqzbxubihoqroxpflqvjuzyjcolauypcjohzocbaswgxyvvzqmrnfxxwcsondrprkrqhhswjxdyrxgngnmbhszcdbxbbecnwafgzfbetbdgkdboethftqsdvsnjjwjocqrbptfkcpznwrvgraqqovdxnwueruvurwhjdqepqkydnpzbpqhhypfexoizrpesmqrsriyjrxpmrcxusirnnuorifjxtpqgxbnfdpvscmbgghcrseiwbfqojfuwhxfqlspedzekwezusjygdgxwzgqguaorvbzuqtxgworthvrkqwkevrtjibearniqivfrkeutwvibacxwwwawuaukbqtlpqegeonpapxieguvttvccjblyflxpop",
		}
		_, err := service.Create(context.Background(), createRequest)
		assert.NotEmpty(t, err)
		assert.ErrorContains(t, err, "description: the length must be no more than 4096")
		assert.ErrorContains(t, err, "state: must be a valid value")
		assert.ErrorContains(t, err, "title: the length must be no more than 128")
	})

	t.Run("Must Create Risk successfully", func(t *testing.T) {
		riskEntity := &entity.Risk{State: "open", Title: "t", Description: "d"}
		createRequest := &risk.CreateRiskRequest{
			State:       "open",
			Title:       "t",
			Description: "d",
		}
		repo.On("Create", mock.Anything, mock.Anything).
			Return(nil).Once()
		repo.On("Get", mock.Anything, mock.AnythingOfType("string")).
			Return(riskEntity, nil).Once()
		r, err := service.Create(context.Background(), createRequest)
		assert.NotEmpty(t, r)
		assert.Empty(t, err)
	})
}
