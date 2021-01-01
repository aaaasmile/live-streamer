<template>
  <v-row justify="center">
    <v-col xs="12" sm="12" md="10" lg="8" xl="6">
      <v-card color="grey lighten-4" flat tile>
        <v-col cols="12">
          <v-row>
            <v-col cols="12" flex-nowrap>
              <v-text-field
                @keydown.enter="enterPress"
                v-model="uriToPlay"
                label="Select an URI"
              ></v-text-field>
              <v-spacer></v-spacer>
              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                  <v-btn
                    icon
                    @click="playUri"
                    :loading="loadingyoutube"
                    v-on="on"
                  >
                    <v-icon>airplay</v-icon>
                  </v-btn>
                </template>
                <span>Stream uri</span>
              </v-tooltip>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-card flat tile>
                <v-card-title>Streaming {{ ListName }}</v-card-title>
                <div class="mx-4">
                  <div class="subtitle-2">URI</div>
                  <div class="subtitle-2 text--secondary">
                    {{ PlayingURI }}
                  </div>
                  <div class="subtitle-2">Title</div>
                  <div class="subtitle-2 text--secondary">
                    {{ PlayingTitle }}
                  </div>
                  <div class="subtitle-2">Description</div>
                  <div class="subtitle-2 text--secondary">
                    {{ PlayingDesc }}
                  </div>
                </div>
              </v-card>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <Playerbar />
            </v-col>
          </v-row>
        </v-col>
      </v-card>
    </v-col>
  </v-row>
</template>