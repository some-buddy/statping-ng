<template>
    <div :style="{ backgroundColor: statusColor }" class="service-status">
      <span v-if="statusActive">
        {{ statusName }}
      </span>
      <span v-else>
        Service is online
      </span>
    </div>
  </template>
  
  <script>
  export default {
    props: {
      statusActive: Boolean, // Indique si le statut intermédiaire est actif
      statusName: String,    // Nom du statut (Minor Outage ou Major Outage)
      statusColor: String,   // Couleur associée au statut
    },
  };
  </script>
  
  <style scoped>
  .service-status {
    padding: 10px;
    border-radius: 5px;
    color: white;
    font-weight: bold;
  }
  </style>
  