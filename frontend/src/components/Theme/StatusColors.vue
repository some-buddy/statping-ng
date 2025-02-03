<template>
    <div :style="cssVariables" class="status-colors-provider">
      <slot />
    </div>
</template>

<script>
export default {
  computed: {
    cssVariables() {
      return {
        '--status-minor-outage': this.$store.state.statusColors.minor_outage,
        '--status-minor-outage-border': this.darken(this.$store.state.statusColors.minor_outage, 0.2),
        '--status-minor-outage-icon': this.darken(this.$store.state.statusColors.minor_outage, 0.3),
        '--status-major-outage': this.$store.state.statusColors.major_outage,
        '--status-major-outage-border': this.darken(this.$store.state.statusColors.major_outage, 0.2),
        '--status-major-outage-icon': this.darken(this.$store.state.statusColors.major_outage, 0.3)
      }
    }
  },

  methods: {
    darken(color, amount) {
      if (!color || !color.startsWith('#')) {
          return '#000000' // Fallback noir
      }
      // Implémentation d'une fonction simple de darken
      return `color-mix(in srgb, ${color}, black ${amount * 100}%)`
    }
  }
}
</script>