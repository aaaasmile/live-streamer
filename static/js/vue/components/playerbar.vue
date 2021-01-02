<template>
  <v-skeleton-loader
    :loading="loadingMeta"
    :transition="transition"
    type="list-item-two-line"
  >
    <v-container>
      <v-row justify="center">
        <v-col xs="12" sm="12" md="10" lg="8" xl="6">
          <v-row justify="center">
            <v-toolbar color="light-green lighten-5" flat>
              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                  <v-btn icon v-on="on" @click="previousTitle">
                    <v-icon>mdi-skip-previous</v-icon>
                  </v-btn>
                </template>
                <span>Previous</span>
              </v-tooltip>

              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                  <v-btn icon v-on="on" @click="nextTitle">
                    <v-icon>mdi-skip-next</v-icon>
                  </v-btn>
                </template>
                <span>Next</span>
              </v-tooltip>

              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                  <v-btn icon v-on="on">
                    <v-icon>mdi-repeat</v-icon>
                  </v-btn>
                </template>
                <span>Repeat</span>
              </v-tooltip>
           
              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                  <v-btn icon v-on="on" @click="syncStatus">
                    <v-icon>mdi-sync</v-icon>
                  </v-btn>
                </template>
                <span>Synchronize status</span>
              </v-tooltip>
             
              <v-btn icon @click="togglePower" :color="ColorPower">
                <v-icon>power_settings_new</v-icon>
              </v-btn>
            </v-toolbar>
          </v-row>
        </v-col>
      </v-row>
    </v-container>
  </v-skeleton-loader>
</template>