<template>
  <v-snackbar v-model="snackbar" :color="colorsnack">
      {{ textMsg }}
      <v-btn dark text @click="clearErrorText">Close</v-btn>
    </v-snackbar>
</template>