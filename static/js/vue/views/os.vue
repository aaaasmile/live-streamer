<template>
  <v-container>
    <v-card color="grey lighten-4" flat tile>
      <v-toolbar flat dense>
        <v-toolbar-title class="subheading grey--text"
          >OS Functions</v-toolbar-title
        >
        <v-spacer></v-spacer>
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <v-btn icon @click="dialogReboot = true" v-on="on">
              <v-icon>rowing</v-icon>
            </v-btn>
          </template>
          <span>Reboot</span>
        </v-tooltip>
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <v-btn icon @click="dialogShutdown = true" v-on="on">
              <v-icon>pan_tool</v-icon>
            </v-btn>
          </template>
          <span>Shutdown</span>
        </v-tooltip>
      </v-toolbar>
      <v-list>
        <v-list-item @click="dialogRestartService = true">
          <v-list-item-icon>
            <v-icon>play_for_work</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Service Restart</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item @click="killAllOmx">
          <v-list-item-icon>
            <v-icon>play_for_work</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Kill all OmxPlayer</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-card>
    <v-container>
      <v-dialog v-model="dialogRestartService" persistent max-width="290">
        <v-card>
          <v-card-title class="headline">Question</v-card-title>
          <v-card-text>Do you want to restart the service?</v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="green darken-1" text @click="restart_service"
              >OK</v-btn
            >
            <v-btn
              color="green darken-1"
              text
              @click="dialogRestartService = false"
              >Cancel</v-btn
            >
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-dialog v-model="dialogShutdown" persistent max-width="290">
        <v-card>
          <v-card-title class="headline">Question</v-card-title>
          <v-card-text>Do you want to shutdown the Pi?</v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="green darken-1" text @click="shutdown">OK</v-btn>
            <v-btn color="green darken-1" text @click="dialogShutdown = false"
              >Cancel</v-btn
            >
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-dialog v-model="dialogReboot" persistent max-width="290">
        <v-card>
          <v-card-title class="headline">Question</v-card-title>
          <v-card-text>Do you want to reboot the Pi?</v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="green darken-1" text @click="rebooot">OK</v-btn>
            <v-btn color="green darken-1" text @click="dialogReboot = false"
              >Cancel</v-btn
            >
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-container>
  </v-container>
</template>