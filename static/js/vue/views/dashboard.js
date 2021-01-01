import API from '../apicaller.js'
import Playerbar from '../components/playerbar.js'

export default {
  components: { Playerbar },
  data() {
    return {
      loadingyoutube: false,
      loadingplaylist: false,
      uriToPlay: '',
    }
  },
  computed: {
    ...Vuex.mapState({
      PlayingURI: state => {
        return state.ps.uri
      },
      PlayingTitle: state => {
        return state.ps.title
      },
      PlayingDesc: state => {
        return state.ps.description
      },
      PlayingInfo: state => {
        return state.ps.info
      },
      PlayingType: state => {
        return state.ps.itemtype
      },
      ListName: state => {
        return state.ps.listname
      },
      PlayingPrev: state => {
        return state.ps.previous
      },
      PlayingNext: state => {
        return state.ps.next
      },
    })
  },
  methods: {
    playUri() {
      if (this.uriToPlay === ''){
        console.log('Nothig to play')
        return
      }
      console.log('call playUri')
      let req = { uri: this.uriToPlay }
      this.loadingyoutube = true
      API.PlayUri(this, req)
    },
    enterPress(){
      console.log('Enter pressed')
      this.playUri()
    }
  },
  template: `
    <v-row justify="center">
      <v-col xs="12" sm="12" md="10" lg="8" xl="6">
        <v-card color="grey lighten-4" flat tile>
          <v-toolbar flat dense>
            <v-toolbar-title class="subheading grey--text"
              >Player Control</v-toolbar-title
            >
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
              <span>Play uri</span>
            </v-tooltip>
          </v-toolbar>
          <v-col cols="12">
            <v-row>
              <v-col cols="12">
                <v-text-field
                  @keydown.enter="enterPress"
                  v-model="uriToPlay"
                  label="Select an URI"
                ></v-text-field>
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12">
                <v-card flat tile>
                  <v-card-title>Playing {{ ListName }}</v-card-title>
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
`
}