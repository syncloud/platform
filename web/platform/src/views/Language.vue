<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>{{ $t('language.title') }}</h1>
        <div class="row-no-gutters settingsblock">
          <div class="col2">
            <div class="setline" style="display: flex; justify-content: center; align-items: center; gap: 10px; padding-top: 20px;">
              <span class="span">{{ $t('language.select') }}:</span>
              <el-select id="settings_language" v-model="locale" @change="onLocaleChange" size="default" style="width: 220px">
                <el-option v-for="l in locales" :key="l.code" :label="l.name" :value="l.code" :id="'settings_lang_' + l.code"/>
              </el-select>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { SUPPORTED_LOCALES, setLocale } from '../i18n'

export default {
  name: 'Language',
  data () {
    return {
      locale: this.$i18n.locale,
      locales: SUPPORTED_LOCALES
    }
  },
  mounted () {
    this.locale = this.$i18n.locale
  },
  methods: {
    onLocaleChange (code) {
      setLocale(code)
    }
  }
}
</script>
<style>
@import '../style/site.css';
</style>
