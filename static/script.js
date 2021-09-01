const spotifyText = document.createTextNode("Spotify");
const bandcampText = document.createTextNode("Bandcamp");
//const tab = document.getElementById("my-list");

async function addCell(id, elem) {
  let tab = document.getElementById(id), tmpLink, artist, album, newRow;
  artist = document.createTextNode(elem.artist);
  album = document.createTextNode(elem.album);
  tmpLink = document.createElement("a");
  tmpLink.appendChild(spotifyText.cloneNode());
  tmpLink.title = "Lien Spotify";
  tmpLink.href = elem.spotifyurl;

  newRow = tab.insertRow(-1);
  newRow.insertCell(0).appendChild(artist);
  newRow.insertCell(1).appendChild(album);
  newRow.insertCell(2).appendChild(tmpLink);

  tmpLink = document.createElement("a");
  switch (id) {
    case "artist-found":
      tmpLink.classList.add("only-artist")
    case "found":
      tmpLink.appendChild(bandcampText.cloneNode());
      tmpLink.title = "Lien Bandcamp";
      tmpLink.href = elem.bandcampurl;
      newRow.insertCell(3).appendChild(tmpLink);
      break;
  }
}

async function test() {
  const data = await fetch('/refresh').then(response => response.json());
  if (data != null) {
    document.getElementById("nb-fait").textContent = data.done;
    document.getElementById("nb-total").textContent = data.todo;

    if (data.albums != null) {
      for (const elem of data.albums) {
        addCell("found", elem);
      }
    }

    if (data.artists != null) {
      for (const elem of data.artists) {
        addCell("artist-found", elem);
      }
    }

    if (data.notfound != null) {
      for (const elem of data.notfound) {
        addCell("notfound", elem);
      }
    }
    /*
    if (data.urls != null) {
      var newRow, artist, album, tmpLink;
      for (const elem of data.urls) {
        artist = document.createTextNode(elem.artist);
        album = document.createTextNode(elem.album);
        tmpLink = document.createElement("a");
        tmpLink.appendChild(spotifyText.cloneNode());
        tmpLink.title = "Lien Spotify";
        tmpLink.href = elem.spotifyurl;

        /*
        let tmp = document.createTextNode(elem);
        let newP = document.createElement("p");
        newP.appendChild(tmp);
        document.getElementById("my-list").appendChild(newP);

        newRow = tab.insertRow(-1);
        newRow.insertCell(0).appendChild(artist);
        newRow.insertCell(1).appendChild(album);
        newRow.insertCell(2).appendChild(tmpLink);

        tmpLink = document.createElement("a");
        tmpLink.appendChild(bandcampText.cloneNode());
        tmpLink.title = "Lien Bandcamp";
        tmpLink.href = elem.bandcampurl;
        newRow.insertCell(3).appendChild(tmpLink);
      }
    }
*/


    if (data.done === data.todo) {
      document.getElementById("inf-loader").remove();
      clearInterval(refreshList);
    }
  }
}

window.onbeforeunload = function() {
  return "Vous perdrez tous les artistes trouvés en rafraichissant la page."
}

const refreshList = setInterval(test, 5000);
