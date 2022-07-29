//
//  ContentView.swift
//  Tangled
//
//  Created by zy on 26/07/22.
//

import SwiftUI
import CoreData

struct ContentView: View {
    init(){
        UITextView.appearance().backgroundColor = .clear
        UITableView.appearance().backgroundColor = .clear
    }
    @Environment(\.managedObjectContext) private var viewContext
    
    @FetchRequest(
        sortDescriptors: [NSSortDescriptor(keyPath: \Item.timestamp, ascending: true)],
        animation: .default)
    private var items: FetchedResults<Item>
    
    @State private var isActive: Bool = false
    @State private var searchText = ""
    @State private var titleText = ""
    @State private var descriptionText = ""
    @FocusState private var isFocused: Bool
    @State var createQuestion: Bool = false
    
    @State var results = [Post]()
    
    var body: some View {
        NavigationView {
            ZStack {
                LinearGradient(gradient: Gradient(colors: [Color(hex: "#FFFCF4"), Color(hex: "#FFF4D9"), Color(hex: "#FFF7E1")]), startPoint: .top, endPoint: .bottom)
                    .edgesIgnoringSafeArea(.vertical)
                VStack {
                    SearchBar(text: $searchText)
                        .frame(
                            minWidth: 0,
                            maxWidth: 350
                        ).padding(.top, 20)
                    ScrollView {
                        ForEach(results.filter({ searchText.isEmpty ? true : $0.title!.contains(searchText)}), id: \.id) { item in
                            NavigationLink(destination: DetailQuestion(item.id!)){
                                VStack(alignment: .leading, spacing: 6) {
                                    Text(item.title ?? "Error Title!")
                                        .font(.subheadline)
                                        .bold()
                                        .foregroundColor(Color(hex: "#383838"))
                                        .padding(.horizontal,24)
                                        .padding(.top, 16)
                                        .multilineTextAlignment(.leading)
                                    Text(item.description ?? "Error Content!")
                                        .font(.caption)
                                        .foregroundColor(Color(hex: "#484848"))
                                        .padding(.horizontal,24)
                                        .padding(.top, 4)
                                        .multilineTextAlignment(.leading)
                                }
                                .frame(
                                    minWidth: 0,
                                    maxWidth: 327,
                                    minHeight: 0,
                                    maxHeight: 124,
                                    alignment: .leading
                                )
                                .clipped()
                                .padding(.bottom)
                                .background(Color.white)
                                .cornerRadius(20)
                                .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                            }
                            .padding(.top, 8)
                            .navigationBarTitleDisplayMode(.inline)
                            // this sets the Back button text when a new screen is pushed
                            .navigationTitle("Questions")
                        }
                        .frame(maxWidth: .infinity)
                    }
                    .padding(.top, 20)
                    .toolbar {
                        ToolbarItem(placement: .principal) {
                            // this sets the screen title in the navigation bar, when the screen is visible
                            Text("")
                        }
                        ToolbarItem(placement: .navigationBarLeading) {
                            Text("Questions")
                                .font(.title2)
                                .bold()
                                .foregroundColor(Color(hex: "#383838"))
                                .padding(.horizontal,16)
                                .multilineTextAlignment(.leading)
                        }
                        ToolbarItem(placement: .navigationBarTrailing) {
                            //                            EditButton()
                            Button(action: {createQuestion = true}) {
                                Label("Add Item", systemImage: "plus")
                            }
                        }
                        //                        ToolbarItem {
                        //                            Button(action: {createQuestion = true}) {
                        //                                Label("Add Item", systemImage: "plus")
                        //                            }
                        //                        }
                    }
                }
                .sheet(isPresented: $createQuestion) {
                    CreateQuestionModal(
                        isPresented: $createQuestion,
                        title: $titleText,
                        description: $descriptionText)
                    .onDisappear(perform: {
                        loadData()
                    })
                }
            }.onAppear(perform: {
                loadData()
            })
        }
    }
    
    private func loadData() {
        guard let url = URL(string: "http://127.0.0.1:8888/post") else {
            print("Invalid URL")
            return
        }
        let request = URLRequest(url: url)
        
        URLSession.shared.dataTask(with: request) { data, response, error in
            if let data = data {
                if let response = try? JSONDecoder().decode([Post].self, from: data) {
                    DispatchQueue.main.async {
                        self.results = response
                    }
                    return
                }
            }
        }.resume()
    }
}

private let itemFormatter: DateFormatter = {
    let formatter = DateFormatter()
    formatter.dateStyle = .short
    formatter.timeStyle = .medium
    return formatter
}()

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView().environment(\.managedObjectContext, PersistenceController.preview.container.viewContext)
    }
}
