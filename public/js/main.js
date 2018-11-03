(function () {
    const application = Stimulus.Application.start();
    const IpsumController = class extends Stimulus.Controller {

        initialize() {
            this.sentences = 21;

            if (document.queryCommandSupported("copy")) {
                this.element.classList.add('clipboard--supported');
            }
        }

        copy() {
            this.body.select();
            document.execCommand('copy');
        }

        get body() {
            return this.targets.find('body');
        }

        set body(value) {
            this.data.set('body', value);
            this.targets.find('body').value = value;
        }

        async generate() {
            console.log(`You're requesting ${this.sentences} sentences`);
            const { data: { body } } = await axios({
                method: 'GET',
                url: `/v1/issa-ipsum?sentences=${this.sentences}`
            })

            this.body = body;
        }
    }

    // register controllers
    application.register('ipsum', IpsumController);
})();