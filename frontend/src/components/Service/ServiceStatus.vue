<script>
export default {
  computed: {
    statusStyle() {
      const baseColor = this.$store.state.services.find(service => service.id === this.service.id).outageColor;
      return {
        backgroundColor: baseColor,
        borderColor: this.darkenColor(baseColor, 0.15)
      }
    }
  }
}
</script>

<template>
  <div class="status-indicator" :style="statusStyle">
    <!-- Icône -->
  </div>
</template>