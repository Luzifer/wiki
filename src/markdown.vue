<template>
  <!-- eslint-disable-next-line vue/no-v-html -->
  <div v-html="render" />
</template>

<script>
import showdown from 'showdown'

export default {
  data() {
    return {
      render: '',
    }
  },

  emits: ['rendered'],

  name: 'WikiMarkdown',

  props: {
    content: {
      default: '',
      type: String,
    },

    prerender: {
      default: null,
      type: Function,
    },
  },

  watch: {
    content(to) {
      let content = to
      if (this.prerender) {
        content = this.prerender(content)
      }

      const converter = new showdown.Converter()
      this.render = converter.makeHtml(content)
      this.$emit('rendered')
    },
  },
}
</script>
