import React from 'react';
import { MDBDataTable } from 'mdbreact';
import { Gitlab } from '@gitbeaker/browser';
import parse from 'docker-parse-image'

const api = new Gitlab({
    host: 'https://gitlab.com',
});

class DatatablePage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      not_found: false,
      done: false,
      data: {
        columns: [
          {
            label: 'Vendor',
            field: 'vendor',
            sort: 'asc',
            width: 100
          },
          {
            label: 'Name',
            field: 'name',
            sort: 'asc',
            width: 270
          },
          {
            label: 'Severity',
            field: 'severity',
            sort: 'asc',
            width: 150
          },
          {
            label: 'Package Name',
            field: 'package_name',
            sort: 'asc',
            width: 270
          },
          {
            label: 'Package Version',
            field: 'package_version',
            sort: 'asc',
            width: 100
          },
          {
            label: 'Solution',
            field: 'solution',
            sort: 'asc',
            width: 270
          }
        ],
        rows: []
      }
    };
  }
  async componentDidUpdate(prevProps) {
    if (this.props.image !== '') {
      let image = this.props.image.toLowerCase()
      let registry = parse(image).registry || 'docker.io'
      let namespace = parse(image).namespace || 'library'
      let image_name = parse(image).repository
      let tag = parse(image).tag || 'latest'
      let json_name = registry + '/' + namespace + '/' + image_name + '/' + tag + '.json'
      this.state.data.rows = []
      let results_json
      if (prevProps.image !== this.props.image) {
        try {
          let result = await api.RepositoryFiles.showRaw(25651036,json_name,'main');
          results_json = JSON.parse(result)['runs'][0]['results']  
          let newRows = []

          for (let i=0; results_json.length > i; i++){
              let item = results_json[i]
              newRows.push({
                vendor: item['vendor'],
                severity: item['severity'],
                name: item['name'],
                package_name: item['package_name'],
                package_version: item['package_version'],
                solution: item['solution']
              })
          }
          this.setState({
            data: {
              ...this.state.data,
              rows: newRows,
            },
            done: true,
            not_found: false
          })
        }
        catch (e) {
          /* eslint eqeqeq: 0 */
          if (e == 'HTTPError: 404'){
            this.setState({
              data: {
                ...this.state.data,
                rows: [],
              },
              done: true,
              not_found: true
            })
          } 
        }
      }
    }
  }

  render() {
    if (this.props.image === '' ){
      return <></>
    }
    if (this.state.done === false){
      return <h5>Finding results....</h5>
    }
    if (this.state.not_found === true){
      return <h5>Image not Found</h5>
    }
    return (
      <MDBDataTable
      responsive={true}
        striped={true}
        bordered={true}
        small={true}
        searching={false}
        theadTextWhite
        theadColor="indigo"
        tbodyColor="white"
        data={this.state.data}
      />
    );
  }
}

export default DatatablePage;