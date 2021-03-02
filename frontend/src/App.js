import React from "react";
import Gallery from "react-photo-gallery";
import axios from "axios";

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      photos: [],
    };
  }

  componentDidMount() {
    axios
      .get(`${process.env.REACT_APP_API_URL}/api/photos`, {
        json: true,
      })
      .then(({ data }) => {
        this.setState({
          photos: data.map((photo) => ({
            src: `${process.env.REACT_APP_MINIO_URL}/thumbs/${photo.filename}`,
            key: photo.id,
            width: photo.width,
            height: photo.height,
          })),
        });
      });
  }

  render() {
    return (
      <div className="App">
        <Gallery photos={this.state.photos} />
      </div>
    );
  }
}

export default App;
