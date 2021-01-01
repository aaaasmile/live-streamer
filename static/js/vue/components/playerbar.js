import API from '../apicaller.js'

export default {
  components: {},
  data() {
    return {
      loadingMeta: false,
      transition: 'scale-transition'
    }
  },
  created() {
    console.log('Request player status')
    let req = {}
    API.GetPlayerState(this, req)
  },
  computed: {
    ...Vuex.mapState({
      Muted: state => {
        return state.ps.mute === "muted"
      },
      Playing: state => {
        return state.ps.player === "playing"
      },
      ColorPower: state => {
        if (state.ps.player !== "off" && state.ps.player !== "undef") {
          return "green"
        } else {
          return "error"
        }
      },
      ColorMute: state => {
        if (state.ps.mute === "muted") {
          return "error"
        } else {
          return "gray"
        }
      },
    }),

  },
  methods: {
    syncStatus() {
      this.loadingMeta = true
      console.log('Sync status')
      let req = {}
      API.GetPlayerState(this, req)
    },
    nextTitle() {
      console.log("Next title")
      let req = {}
      API.NextTitle(this, req)
    },
    previousTitle() {
      console.log("Previous title")
      let req = {}
      API.PreviousTitle(this, req)
    },
    togglePower() {
      if (this.$store.state.ps.player === "on" ||
        this.$store.state.ps.player === "playing" ||
        this.$store.state.ps.player === "pause" ||
        this.$store.state.ps.player === "restart") {
        this.loadingMeta = true
        console.log("Power off")
        let req = { power: "off" }
        API.SetPowerState(this, req)
      } else if (this.$store.state.ps.player === "off" ||
        this.$store.state.ps.player === "undef") {
        this.loadingMeta = true
        console.log("Power on")
        let req = { power: "on" }
        API.SetPowerState(this, req)
      } else {
        console.log("ignore toggle on state ", this.$store.state.ps.player)
      }
    },
    toggleMute() {
      let req = {}
      if (this.$store.state.ps.mute === "muted") {
        console.log('Unmute')
        req.volume = 'unmute'
        API.ChangeVolume(this, req)
      } else {
        console.log('Mute')
        req.volume = 'mute'
        API.ChangeVolume(this, req)
      }
    },
    togglePlayResume() {
      let req = {}
      if (this.$store.state.ps.player === "playing") {
        console.log('Pause URI')
        API.Pause(this, req)
      } else {
        API.Resume(this, req)
      }
    },
    VolumeUp() {
      console.log('Volume Up')
      let req = { volume: 'up' }
      API.ChangeVolume(this, req)
    },
    VolumeDown() {
      console.log('Volume Down')
      let req = { volume: 'down' }
      API.ChangeVolume(this, req)
    }
  },
  template: `
  <v-skeleton-loader
    :loading="loadingMeta"
    :transition="transition"
    type="list-item-two-line"
  >
    <v-container>
      <v-row justify="center">
        <v-col xs="12" sm="12" md="10" lg="8" xl="6">
          <v-row justify="center">
            <v-toolbar flat>
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
                  <v-btn icon v-on="on" @click="togglePlayResume">
                    <v-icon>{{ Playing ? "mdi-pause" : "mdi-play" }}</v-icon>
                  </v-btn>
                </template>
                <span>{{ Playing ? "Pause" : "Play current" }}</span>
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
                    <v-icon>mdi-shuffle</v-icon>
                  </v-btn>
                </template>
                <span>Shuffle</span>
              </v-tooltip>

              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                  <v-btn icon v-on="on">
                    <v-icon>mdi-repeat</v-icon>
                  </v-btn>
                </template>
                <span>Repeat</span>
              </v-tooltip>
            </v-toolbar>
            <v-toolbar flat>
              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                  <v-btn icon v-on="on" @click="syncStatus">
                    <v-icon>mdi-sync</v-icon>
                  </v-btn>
                </template>
                <span>Synchronize status</span>
              </v-tooltip>
              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                  <v-btn icon v-on="on" @click="toggleMute" :color="ColorMute">
                    <v-icon>{{ Muted ? "volume_off" : "volume_mute" }}</v-icon>
                  </v-btn>
                </template>
                <span>{{ Muted ? "Unmute" : "Mute" }}</span>
              </v-tooltip>
              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                  <v-btn icon v-on="on" @click="VolumeDown">
                    <v-icon>volume_down</v-icon>
                  </v-btn>
                </template>
                <span>Volume down</span>
              </v-tooltip>
              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                  <v-btn icon v-on="on" @click="VolumeUp">
                    <v-icon>volume_up</v-icon>
                  </v-btn>
                </template>
                <span>Volume up</span>
              </v-tooltip>
              <v-btn icon @click="togglePower" :color="ColorPower">
                <v-icon>power_settings_new</v-icon>
              </v-btn>
            </v-toolbar>
          </v-row>
        </v-col>
      </v-row>
    </v-container>
  </v-skeleton-loader>`
}
