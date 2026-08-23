import createPreset from 'conventional-changelog-conventionalcommits'

export default createPreset({
  types: [
    {type: 'feat', section: 'Features'},
    {type: 'feature', section: 'Features'},
    {type: 'enhance', section: 'Enhancements'},
    {type: 'fix', section: 'Bug Fixes'},
    {type: 'perf', section: 'Performance Improvements'},
    {type: 'refactor', section: 'Refactoring'},
    {type: 'build', section: 'Build'},
    {type: 'ci', section: 'Continuous Integration'},
    {type: 'docs', section: 'Documentation'},
    {type: 'style', section: 'Styles'},
    {type: 'test', section: 'Tests'},
    {type: 'chore', section: 'Chores'},
    {type: 'revert', section: 'Reverts'}
  ]
})
