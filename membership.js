const { Octokit } = require("@octokit/rest");

const octokit = new Octokit({
    auth: process.env.GITHUB_TOKEN,
});

const membership = async (org, team_slug, username) => {
    try {
        // const { data } = await octokit.teams.getMembershipForUserInOrg({
        //     org,
        //     team_slug: team_slug || 'all',
        //     username,
        // });
        const { data } = octokit.rest.teams.list({
            org,
          });
        console.log(data);
    } catch (error) {
        console.error(error);
    }
}

membership('opentdf', 'cli', 'suchak1');
