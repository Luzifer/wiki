<template>
  <!-- eslint-disable-next-line vue/no-v-html -->
  <div v-html="render" />
</template>

<script>
import showdown from 'showdown'

const classMap = {
  blockquote: 'blockquote',
  table: 'table',
}

const htmlClassBindings = Object.keys(classMap)
  .map(key => ({
    regex: new RegExp(`<${key}(.*)>`, 'g'),
    replace: `<${key} class="${classMap[key]}" $1>`,
    type: 'output',
  }))

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

      const converter = new showdown.Converter({
        extensions: [...htmlClassBindings],
        tables: true,
      })
      converter.setFlavor('github')
      this.render = converter.makeHtml(content)
      this.$emit('rendered')
    },
  },
}
</script>

<style>
.blockquote {
  margin-bottom: 1rem;
  font-size: 1rem;
  border-left: 3px solid rgb(var(--bs-secondary-bg-rgb));
  padding-left: 0.5rem;
}
</style>
