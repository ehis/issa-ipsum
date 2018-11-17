// Vue Application

window.onload = function() {
  const App = new Vue({
    el: "#app",
    data: {
      sentences: "",
      body: ""
    },
    created() {
      this.generate();
    },
    methods: {
      async generate(event) {
        const response = await axios({
          method: "GET",
          url: `/v1/issa-ipsum?sentences=${this.sentences}`
        });

        const {
          data: { body }
        } = response;
        this.body = body;
      }
    }
  });
};
