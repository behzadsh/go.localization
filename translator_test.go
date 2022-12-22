package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var flatDictionary = map[string]map[string]string{
	"user.not_found": {
		"en": "User not found!",
		"zh": "找不到用户!",
		"hi": "उपयोगकर्ता नहीं मिला!",
		"es": "¡Usuario no encontrado!",
		"fr": "Utilisateur non trouvé!",
		"ar": "لم يتم العثور على المستخدم!",
		"fa": "کاربر پیدا نشد!",
	},
	"user.blocked": {
		"en": "The user has been blocked!",
		"zh": "该用户已被屏蔽!",
		"hi": "उपयोगकर्ता को ब्लॉक कर दिया गया है!",
		"es": "La usuario ha sido bloqueada!",
		"fr": "L'utilisateur a été bloqué!",
		"ar": "تم حظر المستخدم!",
		"fa": "کاربر مسدود شده است!",
	},
	"messages.insufficient_balance": {
		"en": "Your wallet balance is insufficient for performing this action.",
		"zh": "您的钱包余额不足以执行此操作。",
		"hi": "इस क्रिया को करने के लिए आपका वॉलेट बैलेंस अपर्याप्त है।",
		"es": "El saldo de su billetera es insuficiente para realizar esta acción.",
		"fr": "Le solde de votre portefeuille est insuffisant pour effectuer cette action.",
		"ar": "رصيد محفظتك غير كافٍ لتنفيذ هذا الإجراء.",
		"fa": "موجودی کیف پول شما برای انجام این عمل کافی نیست.",
	},
	"validation.password_length_is_insufficient": {
		"en": "Password length should be more than 6 characters.",
		"zh": "密码长度应大于 6 个字符。",
		"hi": "पासवर्ड की लंबाई 6 वर्णों से अधिक होनी चाहिए।",
		"es": "La longitud de la contraseña debe tener más de 6 caracteres.",
		"fr": "La longueur du mot de passe doit être supérieure à 6 caractères.",
		"ar": "يجب أن يكون طول كلمة المرور أكبر من 6 أحرف.",
		"fa": "طول رمز عبور باید بیشتر از 6 کاراکتر باشد.",
	},
}

var keyParam = map[string]map[string]string{
	"user.not_found":                             nil,
	"user.blocked":                               nil,
	"messages.insufficient_balance":              nil,
	"validation.password_length_is_insufficient": {"length": "6"},
}

func TestNewTranslator(t *testing.T) {
	t.Run("default-configuration", func(t *testing.T) {
		cfg := DefaultConfigs()
		tr, err := NewTranslator(cfg)

		assert.Nil(t, err)
		assert.NotNil(t, tr)
		assert.Equal(t, cfg, tr.config)
		assert.NotNil(t, tr.dictionaryManager)
	})

	t.Run("empty-configuration", func(t *testing.T) {
		tr, err := NewTranslator(Config{})

		assert.Nil(t, err)
		assert.NotNil(t, tr)
		assert.NotNil(t, tr.dictionaryManager)
		assert.NotEmpty(t, tr.config)
		assert.Equal(t, defaultPath, tr.config.TranslationPath)
		assert.Equal(t, "en", tr.config.DefaultLocale)
		assert.Equal(t, "en", tr.config.FallbackLocale)
	})
}

func TestTranslator_TranslateBy(t *testing.T) {
	tr, _ := NewTranslator(DefaultConfigs())

	t.Run("success", func(t *testing.T) {
		for _, locale := range existingLanguages {
			for key, data := range flatDictionary {
				t.Run(locale, func(t *testing.T) {
					translation := tr.TranslateBy(locale, key, keyParam[key])
					assert.Equal(t, data[locale], translation)
				})
			}
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Run("incorrect-key", func(t *testing.T) {
			key := "incorrect_key"
			translation := tr.TranslateBy(existingLanguages[0], key)

			assert.Equal(t, key, translation)
		})

		t.Run("unknown-topic", func(t *testing.T) {
			key := "unknown_topic.some_key"
			translation := tr.TranslateBy(existingLanguages[0], key)

			assert.Equal(t, key, translation)
		})

		t.Run("unknown-topic-key", func(t *testing.T) {
			key := "user.unknown_key"
			translation := tr.TranslateBy(existingLanguages[0], key)

			assert.Equal(t, key, translation)
		})
	})

	t.Run("fallback", func(t *testing.T) {
		key := "user.not_found"
		translation := tr.TranslateBy("nl", key)

		assert.Equal(t, flatDictionary[key]["en"], translation)
	})
}

func TestTranslator_Translate(t *testing.T) {
	locale := "zh"
	tr, _ := NewTranslator(Config{
		DefaultLocale: locale,
	})

	t.Run("success", func(t *testing.T) {
		for key, data := range flatDictionary {
			translation := tr.Translate(key, keyParam[key])

			assert.Equal(t, data[locale], translation)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Run("incorrect-key", func(t *testing.T) {
			key := "incorrect_key"
			translation := tr.Translate(key)

			assert.Equal(t, key, translation)
		})

		t.Run("unknown-topic", func(t *testing.T) {
			key := "unknown_topic.some_key"
			translation := tr.Translate(key)

			assert.Equal(t, key, translation)
		})

		t.Run("unknown-topic-key", func(t *testing.T) {
			key := "user.unknown_key"
			translation := tr.Translate(key)

			assert.Equal(t, key, translation)
		})
	})

	t.Run("fallback", func(t *testing.T) {
		tr, _ = NewTranslator(Config{
			DefaultLocale:  "nl",
			FallbackLocale: "hi",
		})

		key := "user.not_found"
		translation := tr.Translate(key)

		assert.Equal(t, flatDictionary[key]["hi"], translation)
	})
}
