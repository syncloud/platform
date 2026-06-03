import SButton from './components/ui/SButton.vue'
import SInput from './components/ui/SInput.vue'
import SSwitch from './components/ui/SSwitch.vue'
import SProgress from './components/ui/SProgress.vue'
import SAlert from './components/ui/SAlert.vue'
import SSelect from './components/ui/SSelect.vue'
import SOption from './components/ui/SOption.vue'
import SRadioGroup from './components/ui/SRadioGroup.vue'
import SRadio from './components/ui/SRadio.vue'
import SCheckboxGroup from './components/ui/SCheckboxGroup.vue'
import SCheckbox from './components/ui/SCheckbox.vue'

const components = {
  SButton,
  SInput,
  SSwitch,
  SProgress,
  SAlert,
  SSelect,
  SOption,
  SRadioGroup,
  SRadio,
  SCheckboxGroup,
  SCheckbox
}

export default {
  install (app) {
    for (const name in components) {
      app.component(name, components[name])
    }
  }
}

export { components }
